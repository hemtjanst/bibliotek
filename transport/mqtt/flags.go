package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

// FlagStrFunc is passed to the Flags() function, it is compatible with the standard library
// flag.String() function
type FlagStrFunc func(name, def, usage string) *string

// FlagBoolFunc is passed to the Flags() function, it is compatible with the standard library
// flag.Bool() function
type FlagBoolFunc func(name string, def bool, usage string) *bool

// Flags will use the provided callbacks to set up mqtt-specific cli arguments:
//
//  mqtt.address   (localhost:1883) - Address of server                (MQTT_ADDRESS)
//  mqtt.username                   - Username                         (MQTT_USERNAME)
//  mqtt.password                   - Password                         (MQTT_PASSWORD)
//  mqtt.tls                        - Enable TLS (bool flag)           (MQTT_TLS)
//    mqtt.tls-insecure             - Skip cert validation (bool flag) (MQTT_TLS_INSECURE)
//    mqtt.cn                       - CN of server (i.e. hostname)     (MQTT_COMMON_NAME)
//    mqtt.ca                       - Path to authority certificate    (MQTT_CA_PATH)
//    mqtt.cert                     - Path to client certificate       (MQTT_CERT_PATH)
//    mqtt.key                      - Path to certificate key          (MQTT_KEY_PATH)
//  topic.announce                  - Topic for announcements          (HEMTJANST_TOPIC_ANNOUNCE)
//  topic.discover                  - Topic for discovery              (HEMTJANST_TOPIC_DISCOVER)
//  topic.leave                     - Topic for leaving                (HEMTJANST_TOPIC_LEAVE)
//
// The callback interfaces are compatible with the standard library flag-package.
// If using the flag-package, you can easily set up the mqtt transport with CLI flags
// by doing this
//
// (Note that it's important to call mqtt.Flags() before flag.Parse())
//
//  package main
//  import (
//    "flag"
//    "context"
//    "github.com/hemtjanst/bibliotek/transport/mqtt"
//  )
//  func main() {
//    myCustomFlag := flag.String("custom", "", "Set up your own flags here")
//    mqCfg := mqtt.MustFlags(flag.String, flag.Bool)
//    flag.Parse()
//    transport, err := mqtt.New(context.TODO(), mqCfg())
//  }
func Flags(str FlagStrFunc, b FlagBoolFunc) func() (*Config, error) {
	if str == nil || b == nil {
		return nil
	}

	envAddr := os.Getenv("MQTT_ADDRESS")
	if envAddr == "" {
		envAddr = "localhost:1883"
	}
	envDiscoverTopic := os.Getenv("HEMTJANST_TOPIC_DISCOVER")
	if envDiscoverTopic == "" {
		envDiscoverTopic = "discover"
	}
	envAnnounceTopic := os.Getenv("HEMTJANST_TOPIC_ANNOUNCE")
	if envAnnounceTopic == "" {
		envAnnounceTopic = "announce"
	}
	envLeaveTopic := os.Getenv("HEMTJANST_TOPIC_LEAVE")
	if envLeaveTopic == "" {
		envLeaveTopic = "leave"
	}
	var envPwd string
	if os.Getenv("MQTT_PASSWORD") != "" {
		// Mask environment password in --help
		envPwd = "**********"
	}

	envUsername := os.Getenv("MQTT_USERNAME")
	envCommonName := os.Getenv("MQTT_COMMON_NAME")
	envCaPath := os.Getenv("MQTT_CA_PATH")
	envCertPath := os.Getenv("MQTT_CERT_PATH")
	envKeyPath := os.Getenv("MQTT_KEY_PATH")
	envTls, _ := strconv.ParseBool(os.Getenv("MQTT_TLS"))
	envTlsInsecure, _ := strconv.ParseBool(os.Getenv("MQTT_TLS_INSECURE"))

	flagAddress := str("mqtt.address", envAddr, "Address to MQTT endpoint")
	flagUsername := str("mqtt.username", envUsername, "MQTT Username")
	flagPassword := str("mqtt.password", envPwd, "MQTT Password")
	flagTLS := b("mqtt.tls", envTls, "Enable TLS")
	flagTLSInsecure := b("mqtt.tls-insecure", envTlsInsecure, "Disable TLS certificate validation")
	flagTLSHostname := str("mqtt.cn", envCommonName, "Common name of server certificate (usually the hostname)")
	flagCaPath := str("mqtt.ca", envCaPath, "Path to CA certificate")
	flagCertPath := str("mqtt.cert", envCertPath, "Path to Client certificate")
	flagKeyPath := str("mqtt.key", envKeyPath, "Path to Client certificate key")
	flagAnnounceTopic := str("topic.announce", envAnnounceTopic, "Announce topic for Hemtjänst")
	flagDiscoverTopic := str("topic.discover", envDiscoverTopic, "Discover topic for Hemtjänst")
	flagLeaveTopic := str("topic.leave", envLeaveTopic, "Leave topic for hemtjänst")

	return func() (c *Config, err error) {
		c = &Config{}

		if flagAddress != nil && *flagAddress != "" {
			c.Address = append(c.Address, *flagAddress)
		}
		if flagUsername != nil && *flagUsername != "" {
			c.Username = *flagUsername
		}
		if flagPassword != nil && *flagPassword != "" {
			if envPwd != "" && *flagPassword == envPwd {
				// Revert to using the environment if no other password
				// was set via flags
				c.Password = os.Getenv("MQTT_PASSWORD")
			} else {
				c.Password = *flagPassword
			}
		}

		if flagTLS != nil && *flagTLS {
			skipVerify := false
			certFile := ""
			keyFile := ""
			caFile := ""
			cnName := ""
			var cp *x509.CertPool
			var cert tls.Certificate

			if flagTLSInsecure != nil {
				skipVerify = *flagTLSInsecure
			}
			if flagCertPath != nil {
				certFile = *flagCertPath
			}
			if flagKeyPath != nil {
				keyFile = *flagKeyPath
			}
			if flagCaPath != nil {
				caFile = *flagCaPath
			}
			if flagTLSHostname != nil {
				cnName = *flagTLSHostname
			}

			if caFile != "" {
				b, err := ioutil.ReadFile(caFile)
				if err != nil {
					return nil, err
				}
				cp = x509.NewCertPool()
				if !cp.AppendCertsFromPEM(b) {
					return nil, err
				}
			} else {
				cp, err = x509.SystemCertPool()
				if err != nil {
					return nil, err
				}
			}

			if certFile != "" && keyFile != "" {
				cert, err = tls.LoadX509KeyPair(certFile, keyFile)
				if err != nil {
					return nil, err
				}
			}

			c.TLS = &tls.Config{
				ClientCAs:                cp,
				Certificates:             []tls.Certificate{cert},
				InsecureSkipVerify:       skipVerify,
				MinVersion:               tls.VersionTLS12,
				PreferServerCipherSuites: true,
				ServerName:               cnName,
			}
		}

		if flagAnnounceTopic != nil {
			c.AnnounceTopic = *flagAnnounceTopic
		}
		if flagDiscoverTopic != nil {
			c.DiscoverTopic = *flagDiscoverTopic
		}
		if flagLeaveTopic != nil {
			c.LeaveTopic = *flagLeaveTopic
		}
		return
	}
}

func MustFlags(str FlagStrFunc, b FlagBoolFunc) func() *Config {
	fn := Flags(str, b)
	return func() *Config {
		c, err := fn()
		if err != nil {
			log.Fatal(err)
		}
		return c
	}
}
