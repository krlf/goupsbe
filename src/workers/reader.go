package workers

import (
	"log"
	"sync"
	"time"
	"upsbe/config"
	"upsbe/types"

	"github.com/tarm/serial"
)

func Reader(quitFlag chan bool, workersWg *sync.WaitGroup, upsConfig *config.Config, serialStream chan types.StringStream) {

	workersWg.Add(1)
	defer workersWg.Done()
	log.Print("Reader thread starting...")
	for {
		select {
		case stream := <-serialStream:
			for command := range stream.Write {
				stream.Read <- readSerial(command, upsConfig)
			}
			close(stream.Read)
		case <-quitFlag:
			log.Print("Reader thread exit...")
			return
		}
	}

}

func readSerial(command string, upsConfig *config.Config) string {

	serialConfig := &serial.Config{Name: upsConfig.SerialDeviceGet(), Baud: 38400, ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(serialConfig)
	if err != nil {
		log.Print(err)
		return ""
	}

	_, err = s.Write([]byte(command + "\r\n"))
	if err != nil {
		log.Print(err)
		s.Close()
		return ""
	}

	buf := make([]byte, 512)
	_, err = s.Read(buf)
	if err != nil {
		log.Print(err)
		s.Close()
		return ""
	}

	s.Close()

	return string(buf)
}
