package device

import (
	"fmt"
	"manager_xiaomi/utils"
)

type Bulb struct {
	MiIoDevice
	FwVer     string `json:"fw_ver"`
	MiioVer   string `json:"miio_ver"`
	HwVer     string `json:"hw_ver"`
	WifiFwVer string `json:"wifi_fw_ver"`
	Mac       string `json:"mac"`
	Ssid      string `json:"ssid"`
	Bssid     string `json:"bssid"`
	Rssi      int    `json:"rssi"`
	Primary   int    `json:"primary"`
}

func NewBulb(debug bool, model string, id string, ip string, token []byte) *Bulb {
	bulb := &Bulb{
		MiIoDevice: MiIoDevice{
			deviceModel: model,
			deviceType:  CheckDevice(model),
			Id:          utils.ConvertHex(id),
			Ip:          ip,
			Token:       token,
			debug:       debug,
			request:     1,
		},
	}
	return bulb
}

func (b *Bulb) String() string {
	return fmt.Sprintf(`{%s,"ip":"%s","timestamp":%d}`,
		b.MiIoDevice.Retain(), b.Ip, b.Timestamp)
}

func (b *Bulb) Retain() string {
	return fmt.Sprintf(`{%s}`, b.MiIoDevice.Retain())
}
