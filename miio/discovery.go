package miio

import (
	"net"
	"time"
)

type Discovery struct {
	Reporter chan *Device
}

func NewDiscovery() *Discovery {
	d := &Discovery{Reporter: make(chan *Device)}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: devicePort})
	if err != nil {
		return nil
	}

	go func() {
		// send broadcast hello packet
		for {
			helloPacket, err := NewPacket(HelloPacketDeviceId, nil, uint32(time.Now().Unix()), nil)
			if err == nil {
				hello, err := helloPacket.Pack()
				if err == nil {
					conn.WriteToUDP(hello, &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: devicePort})
				}
			}
			time.Sleep(time.Second * 60)
		}
	}()

	go func() {
		// read answers for hello
		buffer := make([]byte, 0x20)
		for {
			_, sourceAddr, err := conn.ReadFromUDP(buffer)
			if err == nil {
				if packet, err := ParsePacket(0, nil, buffer); err == nil {
					if packet.DeviceId != 0xffffffff {
						device := NewDevice(sourceAddr.IP.String())
						_, err := device.Hello()
						if err != nil {
							continue
						}

						d.Reporter <- device
					}
				}
			}
		}
	}()

	return d
}
