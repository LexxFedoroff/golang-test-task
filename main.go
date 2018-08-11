package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
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

func startDiscovering() {

}

func startMessageLoop(period time.Duration) {
	tickChan := time.Tick(period)
	for range tickChan {
		log.Printf("Send messages to other instance")

		otherPort := 8000
		if otherPort == *port {
			otherPort = 8001
		}

		conn, err := net.Dial("tcp", fmt.Sprintf(":%v", otherPort))
		if err != nil {
			log.Print(err)
			continue
		}

		defer conn.Close()

		fmt.Fprintf(conn, appID.String()+"\n")

		bufReader := bufio.NewReader(conn)
		responce, err := bufReader.ReadString('\n')
		if err != nil {
			log.Print(err)
			return
		}

		log.Print(responce)
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

	initGracefulStop()

	go startDiscovering()
	go startMessageLoop(3 * time.Second)

	listen(handler)
}
