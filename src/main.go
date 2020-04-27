package main

import (
	"./app"
	"./config"
	"./monitor"
	"./reader"
	"./server"
	"./writer"
	"log"
)

func main() {
	log.Print("Hello, 世界")

	config.Read()

	app := &app.App{}
	app.Initialize()

	go reader.Worker(app)
	go writer.Worker(app)
	go monitor.Worker(app)
	go server.Worker(app)

	app.Run()
}
