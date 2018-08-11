package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func (app app) sendMessage(address string) (string, bool) {
	dialer := net.Dialer{Timeout: time.Second}
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		log.Print(err)
		return "", false
	}

	defer conn.Close()

	fmt.Fprintf(conn, app.ID+"\n")

	bufReader := bufio.NewReader(conn)

	response, err := bufReader.ReadString('\n')

	if err != nil {
		log.Print(err)
		return "", false
	}

	return response, true
}

func (app app) startMessageLoop() {
	tickChan := time.Tick(app.MessageLoopPeriod)
	for range tickChan {
		for appInfo := range app.appList.iter() {
			go func(info discoveredApp) {
				response, ok := app.sendMessage(info.Address)
				if ok {
					log.Print(response)
				} else {
					app.appList.remove(info)
					log.Printf("Applicaiton has removed")
				}
			}(appInfo)
		}
	}
}
