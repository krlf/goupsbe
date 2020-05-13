package handler

import (
	"../db"
	"../model"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"strconv"
	"../types"
	"../config"
)

type restConfig struct {
	Managed bool
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

func GetConfig(config *config.Config, w http.ResponseWriter, r *http.Request) {
	cfg := restConfig{
		Managed: config.ChargeManagementEnabledGet() }
	respondJSON(w, http.StatusOK, cfg)
}

func SetConfig(config *config.Config, w http.ResponseWriter, r *http.Request) {
	cfg := restConfig{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cfg); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	config.ChargeManagementEnabledSet(cfg.Managed)

	respondJSON(w, http.StatusOK, cfg)
}
