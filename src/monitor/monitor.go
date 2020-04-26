package monitor

import (
	"log"
	"time"
	"../app"
	"../model"
	"os/exec"
	"../config"
)

func Worker(a *app.App) {

	a.WorkersWg.Add(1)
	defer a.WorkersWg.Done()
	log.Print("Monitor thread starting...")
	tick := time.Tick(config.GetMonitorInterval() * time.Millisecond)
	for {
		select {
		case <-tick:
			a.WriterSerialWrite <- "GET"
			readings := <-a.WriterSerialRead
			v, ok := model.VoltParse(readings)
			if !ok {
				log.Println("[Monitor] Readings are not parsed.")
				continue
			}
			switch {
			case v.Vb1 < 7200:
				a.WriterSerialWrite <- "CHREN"
				<-a.WriterSerialRead
				log.Println("CHREN")
			case v.Vb1 > 7800:
				a.WriterSerialWrite <- "CHRDIS"
				<-a.WriterSerialRead
				log.Println("CHRDIS")
			case v.Vb1 < 5900:
				log.Println("[Monitor] Going to shutdown!")
				cmd := exec.Command("echo", "true", "/shutdown_signal")
				_, err := cmd.Output()
				if err != nil {
					log.Println(err)
				}				
			}
		case <-a.Quit:
			log.Print("Monitor thread exit...")
			return
		}
	}
}