package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

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

func startMessageLoop() {

}

func listen(handler func(net.Conn)) {
	address := ":8000"

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
			log.Fatal(err)
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	log.Printf("connection has been accepted")

	defer func() {
		conn.Close()
		log.Printf("connection has been closed")
	}()

	_, err := io.Copy(conn, conn)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.Println("application is starting...")

	initGracefulStop()

	go startDiscovering()
	go startMessageLoop()

	listen(handler)
}
