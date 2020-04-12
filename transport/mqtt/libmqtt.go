package mqtt

import (
	"fmt"
	"github.com/goiiot/libmqtt"
)

type mqttClient interface {
	// Connect to all specified server with client options
	Connect(handler libmqtt.ConnHandler)

	// Publish a message for the topic
	Publish(packets ...*libmqtt.PublishPacket)

	// Subscribe topic(s)
	Subscribe(topics ...*libmqtt.Topic)

	// UnSubscribe topic(s)
	UnSubscribe(topics ...string)

	// Destroy all client connection
	Destroy(force bool)
}

// Code is a wrapper for mqtt codes with utils for displaying the error.
// The type does not check if the code is an error or not, it always implements the error interface for convenience
type Code byte

// Error returns the same information as String(), with the prefix "mqtt: "
func (c Code) Error() string {
	return "mqtt: " + c.String()
}

// String returns the same information as Text(), with the integer code appended
func (c Code) String() string {
	return fmt.Sprintf("%s (%d)", c.Text(), uint8(c))
}

// Text returns a textual representation of the code
func (c Code) Text() string {
	switch c {
	case libmqtt.CodeSuccess: // 0 - Packet: ConnAck, PubAck, PubRecv, PubRel, PubComp, UnSubAck, Auth
		return "success"
	case libmqtt.CodeGrantedQos1: // 1 - Packet: SubAck
		return "granted QoS 1"
	case libmqtt.CodeGrantedQos2: // 2 - Packet: SubAck
		return "granted QoS 2"
	case libmqtt.CodeDisconnWithWill: // 4 - Packet: DisConn
		return "disconnected with will"
	case libmqtt.CodeNoMatchingSubscribers: // 16 - Packet: PubAck, PubRecv
		return "no matching subscribers"
	case libmqtt.CodeNoSubscriptionExisted: // 17 - Packet: UnSubAck
		return "no subscription existed"
	case libmqtt.CodeContinueAuth: // 24 - Packet: Auth
		return "continue auth"
	case libmqtt.CodeReAuth: // 25 - Packet: Auth
		return "re auth"
	case libmqtt.CodeUnspecifiedError: // 128 - Packet: ConnAck, PubAck, PubRecv, SubAck, UnSubAck, DisConn
		return "unspecified error"
	case libmqtt.CodeMalformedPacket: // 129 - Packet: ConnAck, DisConn
		return "malformed packet"
	case libmqtt.CodeProtoError: // 130 - Packet: ConnAck, DisConn
		return "protocol error"
	case libmqtt.CodeImplementationSpecificError: // 131 - Packet: ConnAck, PubAck, PubRecv, SubAck, UnSubAck, DisConn
		return "implementation specific error"
	case libmqtt.CodeUnsupportedProtoVersion: // 132 - Packet: ConnAck
		return "unsupported protocol version"
	case libmqtt.CodeClientIdNotValid: // 133 - Packet: ConnAck
		return "client id not valid"
	case libmqtt.CodeBadUserPass: // 134 - Packet: ConnAck
		return "bad username or password"
	case libmqtt.CodeNotAuthorized: // 135 - Packet: ConnAck, PubAck, PubRecv, SubAck, UnSubAck, DisConn
		return "not authorized"
	case libmqtt.CodeServerUnavail: // 136 - Packet: ConnAck
		return "server unavailable"
	case libmqtt.CodeServerBusy: // 137 - Packet: ConnAck, DisConn
		return "server busy"
	case libmqtt.CodeBanned: // 138 - Packet: ConnAck
		return "banned"
	case libmqtt.CodeServerShuttingDown: // 139 - Packet: DisConn
		return "server is shutting down"
	case libmqtt.CodeBadAuthenticationMethod: // 140 - Packet: ConnAck, DisConn
		return "bad authentication method"
	case libmqtt.CodeKeepaliveTimeout: // 141 - Packet: DisConn
		return "keepalive timeout"
	case libmqtt.CodeSessionTakenOver: // 142 - Packet: DisConn
		return "session taken over"
	case libmqtt.CodeTopicFilterInvalid: // 143 - Packet: SubAck, UnSubAck, DisConn
		return "topic filter invalid"
	case libmqtt.CodeTopicNameInvalid: // 144 - Packet: ConnAck, PubAck, PubRecv, DisConn
		return "topic name invalid"
	case libmqtt.CodePacketIdentifierInUse: // 145 - Packet: PubAck, PubRecv, PubAck, UnSubAck
		return "packet identifier in use"
	case libmqtt.CodePacketIdentifierNotFound: // 146 - Packet: PubRel, PubComp
		return "packet identifier not found"
	case libmqtt.CodeReceiveMaxExceeded: // 147 - Packet: DisConn
		return "receive max exceeded"
	case libmqtt.CodeTopicAliasInvalid: // 148 - Packet: DisConn
		return "topic alias invalid"
	case libmqtt.CodePacketTooLarge: // 149 - Packet: ConnAck, DisConn
		return "packet too large"
	case libmqtt.CodeMessageRateTooHigh: // 150 - Packet: DisConn
		return "message rate too high"
	case libmqtt.CodeQuotaExceeded: // 151 - Packet: ConnAck, PubAck, PubRec, SubAck, DisConn
		return "quota exceeded"
	case libmqtt.CodeAdministrativeAction: // 152 - Packet: DisConn
		return "administrative action"
	case libmqtt.CodePayloadFormatInvalid: // 153 - Packet: ConnAck, PubAck, PubRecv, DisConn
		return "payload format invalid"
	case libmqtt.CodeRetainNotSupported: // 154 - Packet: ConnAck, DisConn
		return "retain not supported"
	case libmqtt.CodeQosNoSupported: // 155 - Packet: ConnAck, DisConn
		return "QoS not supported"
	case libmqtt.CodeUseAnotherServer: // 156 - Packet: ConnAck, DisConn
		return "use another server"
	case libmqtt.CodeServerMoved: // 157 - Packet: ConnAck, DisConn
		return "server moved"
	case libmqtt.CodeSharedSubscriptionNotSupported: // 158 - Packet: SubAck, DisConn
		return "shared subscription not supported"
	case libmqtt.CodeConnectionRateExceeded: // 159 - Packet: ConnAck, DisConn
		return "connection rate exceeded"
	case libmqtt.CodeMaxConnectTime: // 160 - Packet: DisConn
		return "max connect time reached"
	case libmqtt.CodeSubscriptionIdentifiersNotSupported: // 161 - Packet: SubAck, DisConn
		return "subscription identifiers not supported"
	case libmqtt.CodeWildcardSubscriptionNotSupported: // 162 - Packet: SubAck, DisConn
		return "wildcard subscriptions not supported"
	case 255:
		return "network error"
	default:
		return "unknown"
	}
}
