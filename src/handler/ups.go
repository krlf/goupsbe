package handler

import (
	"../db"
	"../model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"../types"
)

type RestPage struct {
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

	page := RestPage{v, newPg, newSz, records}

	respondJSON(w, http.StatusOK, page)
}
