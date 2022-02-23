package discovery

import (
	"manager_xiaomi/device"
	"manager_xiaomi/miio"
	"net"
	"time"
)

var (
	discoveryPort     = 54321
	discoveryInterval = time.Second * 10
)

func NewDiscovery(debug bool, discovery chan *device.MiIoDevice) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: discoveryPort})
	if err != nil {
		return err
	}

	go func() {
		// send broadcast hello packet
		for {
			helloPacket, err := miio.NewPacket(miio.HelloPacketDeviceId, nil, uint32(time.Now().Unix()), nil)
			if err == nil {
				hello, err := helloPacket.Pack()
				if err == nil {
					conn.WriteToUDP(hello, &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: discoveryPort})
				}
			}
			time.Sleep(discoveryInterval)
		}
	}()

	// read answers for hello
	buffer := make([]byte, 0x20)
	for {
		_, sourceAddr, err := conn.ReadFromUDP(buffer)
		if err == nil {
			if packet, err := miio.ParsePacket(0, nil, buffer); err == nil {
				// looking for device with real deviceId
				if packet.DeviceId != 0xffffffff {
					if device := device.NewMiIoDevice(debug, packet.DeviceId, sourceAddr.IP.String()); device != nil {
						discovery <- device
					}
				}
			}
		}
	}
}
