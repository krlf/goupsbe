package workers

import (
	"../config"
	"../db"
	"../model"
	"log"
	"time"
	"../types"
	"sync"
)

func Writer(quitFlag chan bool, workersWg *sync.WaitGroup, config *config.Config, serialStream chan types.StringStream, db *db.Db) {

	workersWg.Add(1)
	defer workersWg.Done()
	log.Print("Writer thread starting...")
	tick := time.Tick(config.WriterIntervalGet() * time.Millisecond)
	for {
		select {
		case <-tick:
			storeReadings(serialStream, db)
		case <-quitFlag:
			log.Println("Writer thread exit...")
			return
		}
	}

}

func storeReadings(serialStream chan types.StringStream, db *db.Db) {

	stream := types.StringStreamCreate()
	serialStream <- stream;
	defer close(stream.Write)

	stream.Write <- "GET"
	readings := <-stream.Read

	v, ok := model.VoltParse(readings)

	if !ok {
		log.Println("[Writer] Readings are not parsed.")
		return
	}

	db.UpsVoltageInsert(v)

}
