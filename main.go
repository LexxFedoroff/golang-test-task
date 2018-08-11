package main

import (
	"log"
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

func main() {
	log.SetOutput(os.Stdout)
	log.Println("applicaion is starting...")

	initGracefulStop()

	var wait = make(chan int)

	<-wait
}
