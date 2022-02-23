package device

import (
	"fmt"
	"manager_xiaomi/utils"
)

type T struct {
	Magic     string `json:"magic"`
	Len       int    `json:"len"`
	Unknown1  string `json:"unknown1"`
	Id        string `json:"id"`
	Token     string `json:"token"`
	Timestamp string `json:"timestamp"`
	Checksum  string `json:"checksum"`
	Data      []struct {
		Id     int `json:"id"`
		Result struct {
			Life   int `json:"life"`
			Ipflag int `json:"ipflag"`

			Ap struct {
				Ssid  string `json:"ssid"`
				Bssid string `json:"bssid"`
				Rssi  string `json:"rssi"`
				Freq  int    `json:"freq"`
			} `json:"ap"`
			Netif struct {
				LocalIp string `json:"localIp"`
				Mask    string `json:"mask"`
				Gw      string `json:"gw"`
			} `json:"netif"`
			MiioTimes []int `json:"miio_times"`
		} `json:"result"`
		ExeTime int `json:"exe_time"`
	} `json:"data"`
}

type Repeater struct {
	MiIoDevice
	FwVer     string `json:"fw_ver"`
	MiioVer   string `json:"miio_ver"`
	HwVer     string `json:"hw_ver"`
	WifiFwVer string `json:"wifi_fw_ver"`
	ClientVer string `json:"miio_client_ver"`
	Mac       string `json:"mac"`
	Ssid      string `json:"ssid"`
	Bssid     string `json:"bssid"`
	Rssi      int    `json:"rssi"`
}

func NewRepeater(debug bool, model string, id string, ip string, token []byte) *Bulb {
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

func (b *Repeater) String() string {
	return fmt.Sprintf(`{%s,"ip":"%s","id":"%x","token":"%x","timestamp":%d}`,
		b.MiIoDevice.Retain(), b.Ip, b.Id, b.Token, b.Timestamp)
}

func (b *Repeater) Retain() string {
	return fmt.Sprintf(`{%s}`, b.MiIoDevice.Retain())
}

func (r *Repeater) Register(sid string, passwd string, _hidden bool, _explorer bool) error {

	return nil
}
