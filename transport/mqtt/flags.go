package mqtt

import (
	"github.com/goiiot/libmqtt"
)

var (
	flagAddress     *string
	flagUsername    *string
	flagPassword    *string
	flagTLS         *bool
	flagTLSInsecure *bool
	flagTLSHostname *string
	flagCaPath      *string
	flagCertPath    *string
	flagKeyPath     *string
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
func Flags(str FlagStrFunc, b FlagBoolFunc) {
	if str == nil || b == nil {
		return
	}
	flagAddress = str("mqtt.address", "localhost:1883", "Address to MQTT endpoint")
	flagUsername = str("mqtt.username", "", "MQTT Username")
	flagPassword = str("mqtt.password", "", "MQTT Password")
	flagTLS = b("mqtt.tls", false, "Enable TLS")
	flagTLSInsecure = b("mqtt.tls-insecure", false, "Disable TLS certificate validation")
	flagTLSHostname = str("mqtt.cn", "", "Common name of server certificate (usually the hostname)")
	flagCaPath = str("mqtt.ca", "", "Path to CA certificate")
	flagCertPath = str("mqtt.cert", "", "Path to Client certificate")
	flagKeyPath = str("mqtt.key", "", "Path to Client certificate key")
}

func flagOpts() (o []libmqtt.Option, err error) {
	if flagAddress != nil && *flagAddress != "" {
		o = append(o, libmqtt.WithServer(*flagAddress))
	}
	if flagUsername != nil && flagPassword != nil && *flagUsername != "" {
		o = append(o, libmqtt.WithIdentity(*flagUsername, *flagPassword))
	}
	if flagTLS != nil && *flagTLS {
		o = append(o, libmqtt.WithTLS(
			*flagCertPath,
			*flagKeyPath,
			*flagCaPath,
			*flagTLSHostname,
			*flagTLSInsecure,
		))
	}
	return
}
