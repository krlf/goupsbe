package writer

import (
	"../app"
	"../config"
	"../db"
	"../model"
	"log"
	"time"
)

func Worker(a *app.App) {

	a.WorkersWg.Add(1)
	defer a.WorkersWg.Done()
	log.Print("Writer thread starting...")
	tick := time.Tick(config.GetWriterInterval() * time.Millisecond)
	for {
		select {
		case <-tick:
			a.WriterSerialWrite <- "GET"
			readings := <-a.WriterSerialRead
			storeReadings(a.Db, readings)
		case <-a.Quit:
			log.Println("Writer thread exit...")
			return
		}
	}

}

func storeReadings(db *db.Db, readings string) {

	v, ok := model.VoltParse(readings)

	if !ok {
		log.Println("[Writer] Readings are not parsed.")
		return
	}

	db.UpsVoltageInsert(v)

}
