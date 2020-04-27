package reader

import (
	"../app"
	"../config"
	"github.com/tarm/serial"
	"log"
	"time"
)

func Worker(a *app.App) {

	a.WorkersWg.Add(1)
	defer a.WorkersWg.Done()
	log.Print("Reader thread starting...")
	for {
		select {
		case command := <-a.WriterSerialWrite:
			a.WriterSerialRead <- readSerial(command)
		case command := <-a.MonitorSerialWrite:
			a.MonitorSerialRead <- readSerial(command)
		case command := <-a.RestSerialWrite:
			a.RestSerialRead <- readSerial(command)
		case <-a.Quit:
			log.Print("Reader thread exit...")
			return
		}
	}

}

func readSerial(command string) string {

	c := &serial.Config{Name: config.GetSerialDevice(), Baud: 38400, ReadTimeout: time.Millisecond * 500}
	s, err := serial.OpenPort(c)
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
