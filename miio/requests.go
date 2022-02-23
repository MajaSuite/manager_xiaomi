package miio

import (
	"fmt"
)

/*
Valid methods with parameters:
	"miIO.info", nil
	"miIO.get_repeater_sta_info", nil
	"miIO.get_repeater_ap_info", nil
	"miIO.config_router", &DeviceConfiguration{Ssid:"NET",Password:"xxx",Uid:0}
	"miIO.wifi_assoc_state"
	"miIO.switch_wifi_ssid", &WifiConfiguration{Ssid:"NET",Password:"xxx",Hidden:0, WifiExplorer:0}
	"miIO.switch_wifi_explorer", &WifiExplorer{WifiExplorer:0}
	"miIO.get_ota_state"
	"miIO.ota_install"
	"miIO.get_ota_progress"
	"miIO.ota"
	"miIO.xgetKeys"
	"miIO.xdel"
	"miIO.xset"
	"miIO.xget"
	"miIO.stop_diag_mode"
	"miIO.bind_stat"
	"miIO.restore"
	"miIO.get_disable_local_restore"
	"miIO.disable_local_restore"
	"miIO.reboot"
	"miIO.config"
	"miIO.set_xy"
	"get_aging_status"
	"set_ps"
	"get_ps"
*/
type Request struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

func (r Request) String() string {
	return fmt.Sprintf(`{"id":%d,"method":"%s","params":"%s"}`, r.Id, r.Method, r.Params)
}

type Response struct {
	Id     int            `json:"id"`
	Result []string       `json:"result,omitempty"`
	Error  *ResponseError `json:"error,omitempty"`
}

func (r Response) String() string {
	return fmt.Sprintf(`{"id":%d,"result":%v,"error":%s}`, r.Id, r.Result, r.Error.String())
}

type ResponseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r ResponseError) String() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, r.Code, r.Message)
}

// method = "miIO.get_repeater_sta_info"
/* todo
"Access policy: {result.access_policy}"
"Associated stations: {result.associated_stations}",

{ 	"result": { "code": 0,
		"mac": { "sta_5g": { }, "sta_lan": { }, "sta_2g": { } },
		"sta": { "count": 0 }
	},
	"id": 10956,
	"exe_time":610
}
*/

// method = "miIO.get_repeater_ap_info"
/* todo
   "SSID: {result.ssid}\n"
   "Password: {result.password}\n"
   "SSID hidden: {result.ssid_hidden}\n",
*/

// method = "miIO.config_router"
type DeviceConfiguration struct {
	Ssid     string `json:"ssid"`
	Password string `json:"passwd"`
	Uid      int    `json:"uid,omitempty"`
}

func (r DeviceConfiguration) String() string {
	return fmt.Sprintf(`{"ssid":"%s","passwd":"%s","uid":"%s"}`, r.Ssid, r.Password, r.Uid)
}

// method = "miIO.switch_wifi_ssid"
type WifiConfiguration struct {
	Ssid         string `json:"ssid"`
	Password     string `json:"pwd"`
	Hidden       int    `json:"hidden"`
	WifiExplorer int    `json:"wifi_explorer"`
}

func (r WifiConfiguration) String() string {
	return fmt.Sprintf(`{"ssid":"%s","pwd":"%s","hidden":%d,"wifi_explorer":%d}`,
		r.Ssid, r.Password, r.Hidden, r.WifiExplorer)
}

// method = "miIO.switch_wifi_explorer"
type WifiExplorer struct {
	WifiExplorer int `json:"wifi_explorer"`
}

func (r WifiExplorer) String() string {
	return fmt.Sprintf(`{"wifi_explorer":%d}`, r.WifiExplorer)
}
