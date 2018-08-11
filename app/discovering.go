package app

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	multicastAddr   = "239.0.0.0:9999"
	maxDatagramSize = 8192
)

var signature = []byte("INSEcosystem_TestTask\n")

func (app app) listenMulticast() {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	l.SetReadBuffer(maxDatagramSize)
	for {
		buffer := make([]byte, maxDatagramSize)
		n, src, err := l.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		if bytes.Equal(src.IP, app.IP) {
			continue
		}

		if !bytes.Equal(buffer[:n], signature) {
			log.Printf("Invalid signature. Skip instance")
			continue
		}

		var port = 8000 // TODO read port from buffer
		inst := discoveredApp{fmt.Sprintf("%v:%v", src.IP, port)}
		if app.appList.add(inst) {
			log.Printf("New instance app has added: %v", inst.Address)
		}
	}
}

func (app app) pingMulticast() {
	addr, err := net.ResolveUDPAddr("udp", multicastAddr)
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	for range time.Tick(time.Second) {
		conn.Write(signature)
	}
}

func (app app) startDiscovering() {
	go app.listenMulticast()
	go app.pingMulticast()
}
