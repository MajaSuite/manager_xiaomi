package miio

import (
	"encoding/json"
	"fmt"
	"time"
)

type request struct {
	Id     uint32      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

func (r request) String() string {
	return fmt.Sprintf(`{"id":%d,"method":"%s","params":"%s"}`, r.Id, r.Method, r.Params)
}

type response struct {
	Id     int            `json:"id"`
	Result []string       `json:"result,omitempty"`
	Error  *responseError `json:"error,omitempty"`
}

func (r response) String() string {
	return fmt.Sprintf(`{"id":%d,"result":%v,"error":%s}`, r.Id, r.Result, r.Error.String())
}

type responseError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r responseError) String() string {
	return fmt.Sprintf(`{"code":%d,"message":"%s"}`, r.Code, r.Message)
}

type credentials struct {
	SID  string `json:"ssid"`
	Pass string `json:"passwd"`
	Uid  string `json:"uid,omitempty"`
}

func (r credentials) String() string {
	return fmt.Sprintf(`{"ssid":"%s","passwd":"%s","uid":"%s"}`, r.SID, r.Pass, r.Uid)
}

func (x *Miio) Hello() (*Packet, error) {
	hello, err := x.PacketHello()
	if err != nil {
		return nil, err
	}

	// send hello packet
	if _, err := x.SendPacket(hello); err != nil {
		return nil, err
	}

	// receive answer for hello
	recv, err := x.ReceivePacket()
	if err != nil {
		return nil, err
	}

	// parse received packet
	pkt, err := ParsePacket(0, nil, recv)
	if err != nil {
		return nil, err
	}

	x.id = pkt.DeviceId
	x.token = pkt.CheckSum
	x.uptime = pkt.Timestamp
	x.timestamp = uint32(time.Now().Unix())

	return pkt, nil
}

func (x *Miio) Info() (*Packet, error) {
	req := request{
		Id:     x.uptime,
		Method: "miIO.info",
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	p, err := x.Packet(x.id, x.token, uint32(time.Now().Unix())-x.timestamp+x.uptime, payload)
	if err != nil {
		panic(err)
	}

	// send request
	if _, err := x.SendPacket(p); err != nil {
		return nil, err
	}

	// receive answer
	recv, err := x.ReceivePacket()
	if err != nil {
		return nil, err
	}

	// parse received packet
	pkt, err := ParsePacket(x.id, x.token, recv)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

func (x *Miio) Reg(sid string, pass string) (*Packet, error) {
	req := request{
		Id:     x.uptime,
		Method: "miIO.config_router",
		Params: &credentials{SID: sid, Pass: pass},
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	p, err := x.Packet(x.id, x.token, uint32(time.Now().Unix())-x.timestamp+x.uptime, payload)
	if err != nil {
		panic(err)
	}

	// send request
	if _, err := x.SendPacket(p); err != nil {
		return nil, err
	}

	// receive answer
	recv, err := x.ReceivePacket()
	if err != nil {
		return nil, err
	}

	// parse received packet
	pkt, err := ParsePacket(x.id, x.token, recv)
	if err != nil {
		return nil, err
	}

	return pkt, nil
}
