package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"upsbe/config"
	"upsbe/db"
	"upsbe/model"
	"upsbe/types"

	"github.com/gorilla/mux"
)

type restConfig struct {
	Managed              bool
	StartChargingVoltage uint32
	StopChargingVoltage  uint32
	ShutdownVoltage      uint32
}

type restPage struct {
	Content    []model.Volt
	PageNumber int
	PageSize   int
	Records    int
}

func GetVolt(stream types.StringStream, w http.ResponseWriter, r *http.Request) {
	stream.Write <- "GET"
	readings := <-stream.Read
	v, ok := model.VoltParse(readings)
	if ok {
		respondJSON(w, http.StatusOK, v)
	} else {
		respondError(w, http.StatusInternalServerError, "Readings are not parsed.")
	}
}

func GetHist(db *db.Db, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pg := 1
	sz := 10

	if vars["pg"] != "" {
		pg, _ = strconv.Atoi(vars["pg"])
		if pg < 1 {
			pg = 1
		}
	}

	if vars["sz"] != "" {
		sz, _ = strconv.Atoi(vars["sz"])
		if sz < 1 {
			sz = 10
		}
	}

	v, records, newPg, newSz := db.UpsVoltageGet(pg, sz)

	page := restPage{v, newPg, newSz, records}

	respondJSON(w, http.StatusOK, page)
}

func GetConfig(upsConfig *config.Config, w http.ResponseWriter, r *http.Request) {
	cfg := restConfig{
		Managed:              upsConfig.ChargeManagementEnabledGet(),
		StartChargingVoltage: upsConfig.StartChargingVoltageGet(),
		StopChargingVoltage:  upsConfig.StopChargingVoltageGet(),
		ShutdownVoltage:      upsConfig.ShutdownVoltageGet()}
	respondJSON(w, http.StatusOK, cfg)
}

func SetConfig(upsConfig *config.Config, w http.ResponseWriter, r *http.Request) {
	cfg := restConfig{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cfg); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	newConfig := &config.Config{}

	newConfig.ChargeManagementEnabledSet(cfg.Managed)
	newConfig.StartChargingVoltageSet(cfg.StartChargingVoltage)
	newConfig.StopChargingVoltageSet(cfg.StopChargingVoltage)
	newConfig.ShutdownVoltageSet(cfg.ShutdownVoltage)

	if err := upsConfig.Apply(newConfig); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, cfg)
}
