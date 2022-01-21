package miio

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

var (
	miioPort = 54321
)

type Miio struct {
	addr      string
	id        uint32
	token     []byte //16
	timestamp uint32 // time of first packet
	uptime    uint32 // device uptime in seconds
	conn      net.Conn
}

func NewMiio(ip string) *Miio {
	if conn, err := net.Dial("udp4", fmt.Sprintf("%s:%d", ip, miioPort)); err == nil {
		return &Miio{
			addr: ip,
			conn: conn,
		}
	}
	return nil
}

func (x *Miio) Close() error {
	return x.conn.Close()
}

func (x *Miio) PacketHello() ([]byte, error) {
	pkt, err := NewPacket(HelloPacketDeviceId, nil, uint32(time.Now().Unix()), nil)
	if err != nil {
		return nil, err
	}

	return pkt.Pack()
}

func (x *Miio) Packet(deviceId uint32, deviceToken []byte, timestamp uint32, payload []byte) ([]byte, error) {
	pkt, err := NewPacket(deviceId, deviceToken, timestamp, payload)
	if err != nil {
		return nil, err
	}

	return pkt.Pack()
}

func (x *Miio) SendPacket(buf []byte) (int, error) {
	return x.conn.Write(buf)
}

func (x *Miio) ReceivePacket() ([]byte, error) {
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
