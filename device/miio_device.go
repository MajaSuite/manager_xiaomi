package device

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"manager_xiaomi/miio"
	"manager_xiaomi/utils"
	"net"
	"time"
)

var (
	timeout             = time.Second * 3
	devicePort          = 54321
	ErrAlreadyConnected = errors.New("already connected")
)

type MiIoDevice struct {
	deviceModel string   `json:"model"`
	deviceType  Type     `json:"type"`
	Name        string   `json:"name"`
	Token       []byte   `json:"token"`
	VmPeak      int      `json:"VmPeak"`
	VmSize      int      `json:"VmSize"`
	VmFree      int      `json:"VmFree"`
	VmRSS       int      `json:"VmRSS"`
	MemFree     int      `json:"MemFree"`
	Ip          string   `json:"-"`
	Id          uint32   `json:"-"`
	Timestamp   uint32   `json:"-"`
	request     int      `json:"-"`
	conn        net.Conn `json:"-"`
	debug       bool     `json:"-"`
}

func NewMiIoDevice(debug bool, id uint32, ip string) *MiIoDevice {
	return &MiIoDevice{
		debug:   debug,
		Id:      id,
		Ip:      ip,
		request: 1,
	}
}

func (x *MiIoDevice) Type() Type {
	return x.deviceType
}

func (x *MiIoDevice) Model() string {
	return x.deviceModel
}

func (x *MiIoDevice) ID() uint32 {
	return x.Id
}

func (x *MiIoDevice) IP() string {
	return x.Ip
}

func (x *MiIoDevice) Connect(ip string) error {
	if x.Timestamp > 0 {
		return ErrAlreadyConnected
	}

	if ip != "" {
		x.Ip = ip
		x.request = 1
	}

	var err error
	if x.conn, err = net.DialTimeout("udp4", fmt.Sprintf("%s:%d", x.Ip, devicePort), timeout); err != nil {
		return err
	}

	hello, err := x.Hello()
	if err != nil {
		return err
	}

	log.Println("mii connect hello", hello.String())

	// in case of registration - save received id & token
	if x.Id == miio.HelloPacketDeviceId && x.Token == nil {
		log.Println("save registration")
		x.Id = hello.DeviceId
		x.Token = hello.CheckSum
		log.Printf("\ndeviceid %x\ntoken %x\n\n", hello.DeviceId, hello.CheckSum)
	}

	return nil
}

func (x *MiIoDevice) Close() error {
	x.Timestamp = 0
	return x.conn.Close()
}

func (x *MiIoDevice) String() string {
	return fmt.Sprintf(`{%s,"ip":"%s","timestamp":%d}`,
		x.Retain(), x.Ip, x.Timestamp)
}

func (x *MiIoDevice) Retain() string {
	return fmt.Sprintf(`"model":"%s","id":"%x","token":"%x"`, x.deviceModel, x.Id, x.Token)
}

// Hello method should be called before start any communication with device.
func (x *MiIoDevice) Hello() (*miio.Packet, error) {
	helloPacket, err := miio.NewPacket(miio.HelloPacketDeviceId, nil, uint32(time.Now().Unix()), nil)
	if err != nil {
		return nil, err
	}

	hello, err := helloPacket.Pack()
	if err != nil {
		return nil, err
	}

	if _, err := x.SendPacket(hello); err != nil {
		return nil, err
	}

	recv, err := x.ReceivePacket()
	if err != nil {
		return nil, err
	}

	// parse received packet
	pkt, err := miio.ParsePacket(0, nil, recv)
	if err != nil {
		return nil, err
	}

	x.Timestamp = uint32(time.Now().Unix()) - pkt.Timestamp

	return pkt, nil
}

// Prepare packet
func (x *MiIoDevice) Prepare(deviceId uint32, deviceToken []byte, timestamp uint32, payload []byte) ([]byte, error) {
	if x.Timestamp == 0 {
		return nil, fmt.Errorf("device not started with Hello() call")
	}

	if x.Token == nil {
		return nil, fmt.Errorf("token unknown for device. can't decrypt communications")
	}

	pkt, err := miio.NewPacket(deviceId, deviceToken, timestamp, payload)
	if err != nil {
		return nil, err
	}

	return pkt.Pack()
}

// Send packet to prepared connection
func (x *MiIoDevice) Send(method string, params interface{}) (*miio.Packet, error) {
	req := miio.Request{
		Id:     int(uint32(time.Now().Unix()) - x.Timestamp), //x.request,
		Method: method,
		Params: params,
	}
	x.request++

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	log.Println("miio send request:", string(payload))

	// prepare packet
	p, err := x.Prepare(x.Id, x.Token, uint32(time.Now().Unix())-x.Timestamp, payload)
	if err != nil {
		return nil, err
	}

	// send request
	if _, err := x.SendPacket(p); err != nil {
		x.Close()
		return nil, err
	}

	// receive answer
	recv, err := x.ReceivePacket()
	if err != nil {
		x.Close()
		return nil, err
	}

	// parse received packet
	pkt, err := miio.ParsePacket(x.Id, x.Token, recv)
	if err != nil {
		return nil, err
	}

	log.Println("miio send response:", pkt.String())

	return pkt, nil
}

func (x *MiIoDevice) SendPacket(buf []byte) (int, error) {
	x.conn.SetWriteDeadline(time.Now().Add(timeout))
	return x.conn.Write(buf)
}

func (x *MiIoDevice) ReceivePacket() ([]byte, error) {
	var magic, size uint16
	var err error

	x.conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(x.conn)

	header := make([]byte, 4)
	if _, err := reader.Read(header); err != nil {
		return nil, err
	}

	if magic, _, err = utils.ReadInt16(header, 0); err != nil {
		return nil, err
	}
	if magic != 0x2131 {
		return nil, fmt.Errorf("wrong magic %d", magic)
	}

	if size, _, err = utils.ReadInt16(header, 2); err != nil {
		return nil, err
	}

	res := make([]byte, size)
	utils.WriteInt16(res, 0, magic)
	utils.WriteInt16(res, 2, size)
	if _, err = reader.Read(res[4:]); err != nil {
		return nil, err
	}

	return res, nil
}
