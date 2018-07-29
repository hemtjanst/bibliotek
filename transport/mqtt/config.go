package mqtt

import (
	"crypto/tls"
	"errors"
	"github.com/goiiot/libmqtt"
	"github.com/google/uuid"
	"os"
	"path"
	"strings"
	"time"
)

type Config struct {
	// ClientID will be used as mqtt Client ID and LastWillID
	ClientID string

	// Address is a slice of MQTT addresses (host:port / ip:port)
	Address []string

	// Username used to authenticate to MQTT
	Username string

	// Password used to authenticate to MQTT
	Password string

	// TLS configuration
	TLS *tls.Config

	// AnnounceTopic is the prefix announcements will have. Default is "announce"
	AnnounceTopic string

	// DiscoverTopic is where the library will send or listen for discoveries. Default is "discover"
	DiscoverTopic string

	// LeaveTopic is where the will is sent when the client exists. Default is "leave"
	LeaveTopic string

	// DiscoverDelay is the time between first subscribing to announcements, and sending a discover.
	// The delay should be long enough to be able to receive all retained announcements, but not too long
	// to make initial startup too slow. Default is 5 seconds
	DiscoverDelay time.Duration
}

func (c *Config) check() error {
	if c.ClientID == "" {
		var id string
		if len(os.Args) > 0 && len(os.Args[0]) > 0 {
			// Use executable name as first part of id
			id = path.Base(os.Args[0])
		} else {
			id = "htlib"
		}
		id = id + "-" + uuid.New().String()
		c.ClientID = id
	}
	if len(c.Address) == 0 {
		return ErrNoAddress
	}
	for _, v := range []string{c.AnnounceTopic, c.DiscoverTopic, c.LeaveTopic} {
		if strings.ContainsAny(v, "#+") {
			return ErrTopicInvalidChar
		}
	}

	return nil
}

func (c *Config) opts() (o []libmqtt.Option) {
	o = []libmqtt.Option{
		libmqtt.WithServer(c.Address...),
		libmqtt.WithKeepalive(10, 1.2),
		libmqtt.WithLog(libmqtt.Silent),
		libmqtt.WithDialTimeout(5),
		libmqtt.WithClientID(c.ClientID),
		libmqtt.WithWill(c.LeaveTopic, 1, false, []byte(c.ClientID)),
	}

	if c.TLS != nil {
		o = append(o, libmqtt.WithCustomTLS(c.TLS))
	}
	if c.Username != "" {
		o = append(o, libmqtt.WithIdentity(c.Username, c.Password))
	}

	return
}

var (
	ErrNoAddress        = errors.New("at least one address must be provided")
	ErrTopicInvalidChar = errors.New("topic contains invalid characters")
	ErrNoConfig         = errors.New("no config provided")
)
