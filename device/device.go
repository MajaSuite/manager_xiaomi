package device

import (
	"encoding/hex"
)

const (
	NO_TYPE Type = iota
	BULB
	RGB_BULB
	REPEATER
)

type Type byte

func (t Type) String() string {
	switch t {
	case BULB:
		return "Bulb"
	case RGB_BULB:
		return "RGB bulb"
	case REPEATER:
		return "WiFi repeater"
	default:
		return "n/a"
	}
}

type Device interface {
	Type() Type
	Model() string
	ID() uint32
	IP() string
	Connect(ip string) error
	Close() error
	String() string
	Retain() string
}

func CheckDevice(model string) Type {
	switch model {
	//case "yeelink.light.monoa":
	//	return BULB
	//case "yeelink.light.monob":
	//	return BULB
	case "yeelink.light.mono1":
		return BULB
	case "yeelink.light.mono4":
		return BULB
	case "yeelink.light.mono5":
		return BULB
	case "yeelink.light.mono6":
		return BULB
	//"xiaomi.repeater.ccccc":
	//"xiaomi.repeater.qwer":
	//"xiaomi.repeater.ra75":
	//"xiaomi.repeater.rtyui":
	//case "xiaomi.repeater.v1":
	//	return REPEATER
	//case "xiaomi.repeater.v2":
	//	return REPEATER
	//case "xiaomi.repeater.v3":
	//	return REPEATER
	//case "xiaomi.repeater.v6":
	//	return REPEATER
	//case "xiaomi.repeater.v7":
	//	return REPEATER

	default:
		return NO_TYPE
	}
}

/* return nil if device doen't known
 */
func CreateDevice(debug bool, model string, id string, ip string, tokenStr string) Device {
	var dev Device
	token, _ := hex.DecodeString(tokenStr)

	switch CheckDevice(model) {
	case BULB:
		dev = NewBulb(debug, model, id, ip, token)
	case REPEATER:
		dev = NewRepeater(debug, model, id, ip, token)
	}

	return dev
}
