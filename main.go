package main

import (
	"flag"
	"golang-test-task/app"
	"log"
	"os"
)

func main() {
	flag.Parse()

	log.SetOutput(os.Stdout)

	app.NewApp().Run()
}
