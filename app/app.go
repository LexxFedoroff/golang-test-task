package app

import (
	"fmt"
	"golang-test-task/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

// App represents test task application
type App interface {
	Run()
	Stop()
}

type app struct {
	ID                string
	IP                net.IP
	Port              int
	MessageLoopPeriod time.Duration
	appList           discoveredAppList
}

func (app app) Address() string {
	return fmt.Sprintf("%v:%v", app.IP, app.Port)
}

func (app app) Run() {
	log.Printf("Application (%v) is starting...", app.ID)

	app.initGracefulStop()

	go app.startDiscovering()
	go app.startMessageLoop()

	app.listen()
}

func (app app) Stop() {
	log.Printf("Application has stopped")
	os.Exit(0)
}

// NewApp creates new instance of App
func NewApp() App {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return app{
		ID:                u4.String(),
		IP:                *utils.GetOutboundIP(),
		Port:              8000,            // TODO randomize port or pass through arguments
		MessageLoopPeriod: 3 * time.Second, // TODO pass through arguments
		appList:           newDiscoveredAppList(),
	}
}

func (app app) initGracefulStop() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		_ = <-gracefulStop
		app.Stop()
	}()
}
