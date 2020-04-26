package main

import (
	"log"
	"./config"
	"./app"
	"./reader"
	"./monitor"
	"./writer"
	"./server"
)

func main() {
	log.Print("Hello, 世界")

	config.Read()

	app := &app.App{}
	app.Initialize();

	go reader.Worker(app)
	go writer.Worker(app)
	go monitor.Worker(app)
	go server.Worker(app)

	app.Run();
}

