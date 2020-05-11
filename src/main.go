package main

import (
	"./app"
	"./config"
	"./workers"
	"log"
	"./db"
)

func main() {
	log.Print("Hello, 世界")

	config := &config.Config{}
	config.Read()

	db := &db.Db{}
	db.Open(config.DbPathGet())

	app := &app.App{}
	app.Initialize(config, db)

	go workers.Reader(app.QuitFlagGet(), app.WorkersWgGet(), config, app.SerialStreamGet())
	go workers.Writer(app.QuitFlagGet(), app.WorkersWgGet(), config, app.SerialStreamGet(), db)
	go workers.Monitor(app.QuitFlagGet(), app.WorkersWgGet(), config, app.SerialStreamGet())
	go workers.Server(app.QuitFlagGet(), app.WorkersWgGet(), config, app.RouterGet())

	app.Run()

	db.Close()
}
