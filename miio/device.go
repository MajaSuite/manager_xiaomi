package miio

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var (
	devicePort = 54321
)

type Device struct {
	Ip        string   `json:"ip"`
	Id        uint32   `json:"id"`
	Token     []byte   //16
	Timestamp uint32   `json:"-"` // time of first packet
	Uptime    uint32   `json:"-"` // device uptime in seconds
	Model     string   `json:"model"`
	Name      string   `json:"name"`
	conn      net.Conn `json:"-"`
	// vendor
	// type
	// x.deviceVendor, x.deviceType = device.CheckDevice(model)
}

func NewDevice(ip string) *Device {
	if conn, err := net.Dial("udp4", fmt.Sprintf("%s:%d", ip, devicePort)); err == nil {
		return &Device{
			Ip:   ip,
			conn: conn,
		}
	}
	return nil
}

func (d *Device) String() string {
	return fmt.Sprintf(`{"ip":"%s","id":"%x","token":"%x","timestamp":%d,"uptime",%d}`, d.Ip, d.Id, d.Token, d.Timestamp, d.Uptime)
}

func (d *Device) Reconnect() error {

	return nil
}

func (x *Device) Close() error {
	return x.conn.Close()
}

func (x *Device) SendPacket(buf []byte) (int, error) {
	return x.conn.Write(buf)
}

func (x *Device) ReceivePacket() ([]byte, error) {
	var magic, size uint16
	var err error

	reader := bufio.NewReader(x.conn)

	header := make([]byte, 4)
	if _, err := reader.Read(header); err != nil {
		return nil, err
	}

	if magic, _, err = ReadInt16(header, 0); err != nil {
		return nil, err
	}
	if magic != 0x2131 {
		return nil, fmt.Errorf("wrong magic %d", magic)
	}

	if size, _, err = ReadInt16(header, 2); err != nil {
		return nil, err
	}

	res := make([]byte, size)
	WriteInt16(res, 0, magic)
	WriteInt16(res, 2, size)
	if _, err = reader.Read(res[4:]); err != nil {
		return nil, err
	}

	return res, nil
}

func (x *Device) PacketHello() ([]byte, error) {
	pkt, err := NewPacket(HelloPacketDeviceId, nil, uint32(time.Now().Unix()), nil)
	if err != nil {
		return nil, err
	}

	return pkt.Pack()
}

func (x *Device) Packet(deviceId uint32, deviceToken []byte, timestamp uint32, payload []byte) ([]byte, error) {
	if x.Uptime == 0 || x.Timestamp == 0 {
		return nil, fmt.Errorf("device not started with Hello() call")
	}

	if x.Token == nil {
		return nil, fmt.Errorf("token unknown for device. can't decrypt communications")
	}

	pkt, err := NewPacket(deviceId, deviceToken, timestamp, payload)
	if err != nil {
		return nil, err
	}

	return pkt.Pack()
}
