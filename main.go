package main

import (
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/client"
	"github.com/MajaSuite/mqtt/packet"
	"log"
	"manager_xiaomi/device"
	"manager_xiaomi/discovery"
	"manager_xiaomi/miio"
)

var (
	debug     = flag.Bool("debug", false, "print debuging hex dumps")
	srv       = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid  = flag.String("clientid", "xiaomi-1", "client id for mqtt server")
	keepalive = flag.Int("keepalive", 30, "keepalive timeout for mqtt server")
	login     = flag.String("login", "", "login string for mqtt server")
	pass      = flag.String("pass", "", "password string for mqtt server")
	qos       = flag.Int("qos", 0, "qos to send/receive from mqtt")
	reg       = flag.Bool("reg", false, "to register new device")
	sid       = flag.String("sid", "myhome", "network name for registration")
	key       = flag.String("key", "mypass", "network key for registration")
	ip        = flag.String("ip", "192.168.1.1", "ip address of new device")
	uid       = flag.Int("uid", 0, "mihome uid")
)

func main() {
	flag.Parse()

	if *reg {
		log.Println("new device registration")

		device := device.NewMiIoDevice(*debug, miio.HelloPacketDeviceId, *ip)
		err := device.Connect(*ip)
		if err != nil {
			log.Println("error connect", err)
			return
		}

		_, err = device.Send("miIO.config_router",
			&miio.DeviceConfiguration{Ssid: *sid, Password: *key, Uid: *uid})
		if err != nil {
			panic(err)
		}

		return
	}

	log.Println("starting manager_xiaomi")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttId uint16 = 1
	mqtt, err := client.Connect(*srv, *clientid, uint16(*keepalive), false, *login, *pass /* *debug */, false)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttId
	sp.Topics = []packet.SubscribePayload{{Topic: "xiaomi/#", QoS: 1}}
	mqtt.Send <- sp
	mqttId++

	devices := make(map[uint32]device.Device)

	log.Println("start xiaomi discovery")
	d := make(chan *device.MiIoDevice)
	go discovery.NewDiscovery(*debug, d)

	for {
		select {
		case pkt := <-mqtt.Receive:
			// todo
			log.Println("mqtt receive", pkt)

		case dev := <-d:
			if dev == nil {
				continue
			}
			if devices[dev.ID()] != nil && devices[dev.ID()].IP() == "" {
				log.Println("device", devices[dev.ID()])
				err := devices[dev.ID()].Connect(dev.Ip)
				if err != nil {
					log.Println("error connect:", err)
				} else {
					p := packet.NewPublish()
					p.Id = mqttId
					p.Topic = fmt.Sprintf("xiaomi/%x", dev.ID())
					p.QoS = packet.QoS(*qos)
					mqttId++
					p.Payload = devices[dev.ID()].String()
					log.Println("payload=", p.Payload)
					mqtt.Send <- p
				}
			}
		}
	}
}
