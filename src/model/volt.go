package model

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type Volt struct {
	Vcc     uint32
	Vin     uint32
	Vcin    uint32
	Vout    uint32
	Vb1     uint32
	Vb2     uint32
	Created time.Time
}

func VoltParse(values string) (*Volt, bool) {

	parsed := false

	if values == "" {
		log.Println("[Parse] Empty string.")
		return nil, parsed
	}

	result := Volt{}

	lines := strings.Split(values, "\n")

	keyCounter := 0

	for _, kv := range lines {
		kvs := strings.Split(kv, ":")
		if len(kvs) == 2 {
			key := strings.Trim(kvs[0], " \r\n")
			value := strings.Trim(kvs[1], " \r\n")
			switch key {
			case "VCC":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vcc = uint32(u)
					keyCounter++
				}
			case "IN":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vin = uint32(u)
					keyCounter++
				}
			case "CIN":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vcin = uint32(u)
					keyCounter++
				}
			case "OUT":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vout = uint32(u)
					keyCounter++
				}
			case "B1":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vb1 = uint32(u)
					keyCounter++
				}
			case "B2":
				u, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					result.Vb2 = uint32(u)
					keyCounter++
				}
			}
		}
	}

	if keyCounter == 6 {
		parsed = true
		result.Created = time.Now()
	}

	return &result, parsed
}
