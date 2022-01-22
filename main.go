package main

import (
	"flag"
	"fmt"
	"github.com/MajaSuite/mqtt/packet"
	"github.com/MajaSuite/mqtt/transport"
	"log"
	"manager_xiaomi/miio"
	"time"
)

var (
	srv       = flag.String("mqtt", "127.0.0.1:1883", "mqtt server address")
	clientid  = flag.String("clientid", "xiaomi-1", "client id for mqtt server")
	keepalive = flag.Int("keepalive", 30, "keepalive timeout for mqtt server")
	login     = flag.String("login", "", "login string for mqtt server")
	pass      = flag.String("pass", "", "password string for mqtt server")
	debug     = flag.Bool("debug", false, "print debuging hex dumps")
	reg       = flag.Bool("reg", false, "to register new device")
	sid       = flag.String("sid", "myhome", "network name for registration")
	key       = flag.String("key", "mypass", "network key for registration")
)

func main() {
	flag.Parse()

	if *reg {
		log.Println("new device registration")

		device := miio.NewDevice("192.168.4.1")
		if device == nil {
			panic("can't init device")
		}

		hello, err := device.Hello()
		if err != nil {
			panic(err)
		}
		log.Println("Hello packet", hello.String())

		info, err := device.Info()
		log.Println("Info packet", info.String())

		register, err := device.Reg(*sid, *key)
		log.Println("Register packet", register.String())
		return
	}

	log.Println("starting manager_xiaomi ...")

	// connect to mqtt
	log.Println("try connect to mqtt")
	var mqttId uint16 = 1
	mqtt, err := transport.Connect(*srv, *clientid, uint16(*keepalive), *login, *pass, *debug)
	if err != nil {
		panic("can't connect to mqtt server " + err.Error())
	}
	go mqtt.Start()

	log.Println("subscribe to managed topics")
	sp := packet.NewSubscribe()
	sp.Id = mqttId
	sp.Topics = []packet.SubscribePayload{{Topic: "xiaomi/#", QoS: 1}}
	mqtt.Sendout <- sp
	mqttId++

	devices := make(map[uint32]*miio.Device)
	// fetch command data from mqtt server
	go func() {
		//
	}()
	time.Sleep(time.Second * 2)

	discovery := miio.NewDiscovery()
	for d := range discovery.Reporter {
		if devices[d.Id] == nil {
			log.Printf("new device: %s", d.String())
			//d.Start(updates)

			p := packet.NewPublish()
			p.Id = mqttId
			p.Topic = fmt.Sprintf("xiaomi/%s", d.Id)
			p.QoS = 1
			p.Payload = fmt.Sprintf(`{"id":"%s","model":"%s","name":"%s"}`, d.Id, d.Model, d.Name)
			p.Retain = true
			//mqtt.Sendout <- p
			mqttId++
		} else {
			d.Close()
			log.Printf("device: %s", d.String())

		}
	}
}
