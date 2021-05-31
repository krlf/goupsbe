package workers

import (
	"log"
	"os"
	"sync"
	"time"

	"upsbe/config"
	"upsbe/model"
	"upsbe/types"
)

func Monitor(quitFlag chan bool, workersWg *sync.WaitGroup, upsConfig *config.Config, serialStream chan types.StringStream) {

	workersWg.Add(1)
	defer workersWg.Done()
	log.Print("Monitor thread starting...")
	var chargingCounter uint8
	var chargingState bool
	tick := time.Tick(upsConfig.MonitorIntervalGet() * time.Millisecond)
	for {
		select {
		case <-upsConfig.ConfigUpdatedTriggerGet():
			log.Print("[Monitor] Config updated.")
			manageVoltageData(serialStream, upsConfig, 1, &chargingState)
		case <-tick:
			if chargingCounter == 10 {
				chargingCounter = 0
			} else {
				chargingCounter++
			}
			manageVoltageData(serialStream, upsConfig, chargingCounter, &chargingState)
		case <-quitFlag:
			log.Print("Monitor thread exit...")
			return
		}
	}
}

func manageVoltageData(serialStream chan types.StringStream, upsConfig *config.Config, chargingCounter uint8, chargingState *bool) {

	stream := types.StringStreamCreate()
	serialStream <- stream
	defer close(stream.Write)

	if !upsConfig.ChargeManagementEnabledGet() {
		stream.Write <- "CHREN"
		<-stream.Read
		return
	}

	if chargingCounter == 2 && *chargingState {
		stream.Write <- "CHREN"
		<-stream.Read
		log.Println("CHREN after check")
	}

	stream.Write <- "GET"
	readings := <-stream.Read
	v, ok := model.VoltParse(readings)
	if !ok {
		log.Println("[Monitor] Readings are not parsed.")
		return
	}

	switch {
	case v.Vb1 < upsConfig.StartChargingVoltageGet():
		*chargingState = true
		stream.Write <- "CHREN"
		<-stream.Read
		log.Println("CHREN")
	case v.Vb1 > upsConfig.StopChargingVoltageGet():
		switch chargingCounter {
		case 0:
			stream.Write <- "CHRDIS"
			<-stream.Read
			log.Println("CHRDIS for check")
		case 1:
			*chargingState = false
			stream.Write <- "CHRDIS"
			<-stream.Read
			log.Println("CHRDIS")
		}
	case v.Vb1 < upsConfig.ShutdownVoltageGet():
		log.Println("[Monitor] Going to shutdown!")
		f, err := os.OpenFile("/shutdown_signal/flag", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err = f.WriteString("true\n")
		if err != nil {
			log.Fatal(err)
		}
	}
}
