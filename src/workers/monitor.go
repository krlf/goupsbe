package workers

import (
	"../config"
	"../model"
	"log"
	"os/exec"
	"time"
	"../types"
	"sync"
)

func Monitor(quitFlag chan bool, workersWg *sync.WaitGroup, config *config.Config, serialStream chan types.StringStream) {

	workersWg.Add(1)
	defer workersWg.Done()
	log.Print("Monitor thread starting...")
	tick := time.Tick(config.MonitorIntervalGet() * time.Millisecond)
	for {
		select {
		case <-config.ConfigUpdatedTriggerGet():
			log.Print("[Monitor] Config updated.")
			manageVoltageData(serialStream, config)
		case <-tick:
			manageVoltageData(serialStream, config)
		case <-quitFlag:
			log.Print("Monitor thread exit...")
			return
		}
	}
}

func manageVoltageData(serialStream chan types.StringStream, config *config.Config) {

	stream := types.StringStreamCreate()
	serialStream <- stream;
	defer close(stream.Write)

	if !config.ChargeManagementEnabledGet() {
		stream.Write <- "CHREN"
		<-stream.Read
		return
	}

	stream.Write <- "GET"
	readings := <-stream.Read
	v, ok := model.VoltParse(readings)
	if !ok {
		log.Println("[Monitor] Readings are not parsed.")
		return
	}
	switch {
	case v.Vb1 < 7200:
		stream.Write <- "CHREN"
		<-stream.Read
		log.Println("CHREN")
	case v.Vb1 > 7800:
		stream.Write <- "CHRDIS"
		<-stream.Read
		log.Println("CHRDIS")
	case v.Vb1 < 5900:
		log.Println("[Monitor] Going to shutdown!")
		cmd := exec.Command("echo", "true", "/shutdown_signal")
		_, err := cmd.Output()
		if err != nil {
			log.Println(err)
		}
	}
}