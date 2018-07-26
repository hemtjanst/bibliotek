package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

// FlagStrFunc is passed to the Flags() function, it is compatible with the standard library
// flag.String() function
type FlagStrFunc func(name, def, usage string) *string

// FlagBoolFunc is passed to the Flags() function, it is compatible with the standard library
// flag.Bool() function
type FlagBoolFunc func(name string, def bool, usage string) *bool

// Flags will use the provided callbacks to set up mqtt-specific cli arguments:
//
//  mqtt.address   (localhost:1883) - Address of server
//  mqtt.username                   - Username
//  mqtt.password                   - Password
//  mqtt.tls                        - Enable TLS (bool flag)
//    mqtt.tls-insecure             - Skip cert validation (bool flag)
//    mqtt.cn                       - CN of server (i.e. hostname)
//    mqtt.ca                       - Path to authority certificate
//    mqtt.cert                     - Path to client certificate
//    mqtt.key                      - Path to certificate key
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
//    mqtt.Flags(flag.String, flag.Bool)
//    flag.Parse()
//    transport, err := mqtt.New(context.TODO(), "")
//  }
func Flags(str FlagStrFunc, b FlagBoolFunc) func() (*Config, error) {
	if str == nil || b == nil {
		return nil
	}
	flagAddress := str("mqtt.address", "localhost:1883", "Address to MQTT endpoint")
	flagUsername := str("mqtt.username", "", "MQTT Username")
	flagPassword := str("mqtt.password", "", "MQTT Password")
	flagTLS := b("mqtt.tls", false, "Enable TLS")
	flagTLSInsecure := b("mqtt.tls-insecure", false, "Disable TLS certificate validation")
	flagTLSHostname := str("mqtt.cn", "", "Common name of server certificate (usually the hostname)")
	flagCaPath := str("mqtt.ca", "", "Path to CA certificate")
	flagCertPath := str("mqtt.cert", "", "Path to Client certificate")
	flagKeyPath := str("mqtt.key", "", "Path to Client certificate key")
	flagAnnounceTopic := str("topic.announce", "announce", "Announce topic for Hemtjänst")
	flagDiscoverTopic := str("topic.discover", "discover", "Discover topic for Hemtjänst")
	flagLeaveTopic := str("topic.leave", "leave", "Leave topic for hemtjänst")

	return func() (c *Config, err error) {
		c = &Config{}

		if flagAddress != nil && *flagAddress != "" {
			c.Address = append(c.Address, *flagAddress)
		}
		if flagUsername != nil && *flagUsername != "" {
			c.Username = *flagUsername
		}
		if flagPassword != nil && *flagPassword != "" {
			c.Password = *flagPassword
		}
		if flagTLS != nil && *flagTLS {
			skipVerify := false
			certFile := ""
			keyFile := ""
			caFile := ""
			cnName := ""

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

			b, err := ioutil.ReadFile(caFile)
			if err != nil {
				return nil, err
			}
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM(b) {
				return nil, err
			}
			cert, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				return nil, err
			}

			c.TLS = &tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: skipVerify,
				ClientCAs:          cp,
				ServerName:         cnName,
			}
		}
		if flagAnnounceTopic != nil && *flagAnnounceTopic != "" {
			c.AnnounceTopic = *flagAnnounceTopic
		}
		if flagDiscoverTopic != nil && *flagDiscoverTopic != "" {
			c.DiscoverTopic = *flagDiscoverTopic
		}
		if flagLeaveTopic != nil && *flagLeaveTopic != "" {
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
