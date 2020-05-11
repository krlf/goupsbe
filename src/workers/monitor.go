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
	/*
	configUpdateTriggerStream := config.ConfigUpdatedTriggerGet()
	reconfigurtionTrigger := <- configUpdateTriggerStream
	*/
	for {
		select {
		/*
		case <-reconfigurtionTrigger.Flag:
			log.Print("[Monitor] Config updated. It is needed to re-read parameter.")
			reconfigurtionTrigger = <- configUpdateTriggerStream
		*/
		case <-tick:
			manageVoltageData(serialStream)
		case <-quitFlag:
			log.Print("Monitor thread exit...")
			return
		}
	}
}

func manageVoltageData(serialStream chan types.StringStream) {

	stream := types.StringStreamCreate()
	serialStream <- stream;
	defer close(stream.Write)

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