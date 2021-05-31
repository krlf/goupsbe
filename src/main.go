package main

import (
	"log"
	"upsbe/app"
	"upsbe/config"
	"upsbe/db"
	"upsbe/workers"
)

func main() {
	log.Print("Hello, 世界")

	upsConfig := &config.Config{}
	upsConfig.Read()

	db := &db.Db{}
	db.Open(upsConfig.DbPathGet())

	app := &app.App{}
	app.Initialize(upsConfig, db)

	go workers.Reader(app.QuitFlagGet(), app.WorkersWgGet(), upsConfig, app.SerialStreamGet())
	go workers.Writer(app.QuitFlagGet(), app.WorkersWgGet(), upsConfig, app.SerialStreamGet(), db)
	go workers.Monitor(app.QuitFlagGet(), app.WorkersWgGet(), upsConfig, app.SerialStreamGet())
	go workers.Server(app.QuitFlagGet(), app.WorkersWgGet(), upsConfig, app.RouterGet())

	app.Run()

	db.Close()
}
