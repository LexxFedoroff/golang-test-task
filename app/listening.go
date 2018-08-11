package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func (app app) listen() {
	address := app.Address()
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Listening has been started on address `%v`", address)
	defer func() {
		l.Close()
		log.Printf("Listening has been stopped on address `%v`", address)
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}

		go app.handler(conn)
	}
}

func (app app) handler(conn net.Conn) {
	defer conn.Close()

	bufReader := bufio.NewReader(conn)
	_, err := bufReader.ReadString('\n')
	if err != nil {
		log.Print(err)
		return
	}

	_, err = fmt.Fprintf(conn, fmt.Sprintf("Hello from %v\n", app.ID))
	if err != nil {
		log.Print(err)
		return
	}
}
