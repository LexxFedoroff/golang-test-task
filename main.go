package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nu7hatch/gouuid"
)

func generateUUID() *uuid.UUID {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return u4
}

var appID = generateUUID()

var port = flag.Int("port", 8000, "")

func initGracefulStop() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		_ = <-gracefulStop
		log.Printf("application has stopped")
		os.Exit(0)
	}()
}

const (
	srvAddr         = "239.0.0.0:9999"
	maxDatagramSize = 8192
)

func getOutboundIP() *net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return &localAddr.IP
}

var appIP = getOutboundIP()

func startDiscovering() {
	go func() {
		addr, err := net.ResolveUDPAddr("udp", srvAddr)
		if err != nil {
			log.Fatal(err)
		}
		l, err := net.ListenMulticastUDP("udp", nil, addr)
		if err != nil {
			log.Fatal(err)
		}
		l.SetReadBuffer(maxDatagramSize)
		for {
			b := make([]byte, maxDatagramSize)
			n, src, err := l.ReadFromUDP(b)
			if err != nil {
				log.Fatal("ReadFromUDP failed:", err)
			}

			if !bytes.Equal(src.IP, *appIP) {
				log.Println(n, "bytes read from", src)
				log.Println(hex.Dump(b[:n]))
			}
		}
	}()

	go func() {
		addr, err := net.ResolveUDPAddr("udp", srvAddr)
		if err != nil {
			log.Fatal(err)
		}
		c, err := net.DialUDP("udp", nil, addr)
		for {
			c.Write([]byte("hello, world\n"))
			time.Sleep(3 * time.Second)
		}
	}()
}

func sendMessage(address string) *string {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Print(err)
		return nil
	}

	defer conn.Close()

	fmt.Fprintf(conn, appID.String()+"\n")

	bufReader := bufio.NewReader(conn)
	response, err := bufReader.ReadString('\n')
	if err != nil {
		log.Print(err)
		return nil
	}

	return &response
}

type appInstance struct {
	Address string
}

type appInstances struct {
	mut   sync.RWMutex
	items []appInstance
}

var appList = appInstances{}

func (list *appInstances) append(app appInstance) {
	list.mut.Lock()
	defer list.mut.Unlock()

	list.items = append(list.items, app)
}

func getApp() appInstance {
	otherPort := 8000
	if otherPort == *port {
		otherPort = 8001
	}

	addr := fmt.Sprintf(":%v", otherPort)

	return appInstance{addr}
}

func (list *appInstances) iter() <-chan appInstance {
	c := make(chan appInstance)

	iter := func() {
		list.mut.RLock()
		defer list.mut.RUnlock()

		for _, app := range list.items {
			c <- app
		}

		close(c)
	}

	go iter()

	return c
}

func startMessageLoop(period time.Duration) {
	tickChan := time.Tick(period)
	for range tickChan {
		log.Printf("Send messages to other instance")

		for inst := range appList.iter() {
			response := sendMessage(inst.Address)
			if response != nil {
				log.Print(*response)
			}
		}
	}
}

func listen(handler func(net.Conn)) {
	address := fmt.Sprintf(":%v", *port)

	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening has been started on address `%v`", address)
	defer func() {
		l.Close()
		log.Printf("listening has been stopped on address `%v`", address)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	log.Print("connection has been accepted")

	defer func() {
		conn.Close()
		log.Print("connection has been closed")
	}()

	bufReader := bufio.NewReader(conn)
	request, err := bufReader.ReadString('\n')
	if err != nil {
		log.Print(err)
		return
	}

	_, err = fmt.Fprintf(conn, fmt.Sprintf("Hello, %v\n", request))
	if err != nil {
		log.Print(err)
		return
	}
}

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.Print("application is starting...")
	log.Printf("Application UUID %v", appID)
	log.Printf("Application IP address %v", appIP)

	initGracefulStop()

	go startDiscovering()
	go startMessageLoop(3 * time.Second)

	listen(handler)
}
