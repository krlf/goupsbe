package db

import (
	"../config"
	"../model"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

type Db struct {
	DbConn *sql.DB
}

func (db *Db) Open() {
	dbConn, err := sql.Open("sqlite3", config.GetDbPath())
	if err != nil {
		log.Fatal(err)
	}
	db.DbConn = dbConn

	stmt, err := dbConn.Prepare(
		"CREATE TABLE IF NOT EXISTS ups " +
			"(uid INTEGER PRIMARY KEY AUTOINCREMENT, " +
			"vcc INTEGER, " +
			"vin INTEGER, " +
			"vcin INTEGER, " +
			"vout INTEGER, " +
			"vb1 INTEGER, " +
			"vb2 INTEGER, " +
			"created TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("DB connection ready.")
}

func (db *Db) Close() {
	if db.DbConn != nil {
		db.DbConn.Close()
		log.Print("DB connection closing.")
	}
}

func (db *Db) UpsVoltageInsert(v *model.Volt) {

	if db.DbConn == nil {
		log.Print("DB connection is not ready.")
		return
	}

	stmt, err := db.DbConn.Prepare("INSERT INTO ups(vcc, vin, vcin, vout, vb1, vb2) values(?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = stmt.Exec(v.Vcc, v.Vin, v.Vcin, v.Vout, v.Vb1, v.Vb2)
	if err != nil {
		log.Println(err)
		return
	}

}

func (db *Db) UpsVoltageGet(pg int, sz int) (v []model.Volt, records int, newPg int, newSz int) {

	rows, err := db.DbConn.Query("SELECT COUNT(*) as count FROM ups")
	if err != nil {
		log.Println("[UpsVoltageGet] Query error:")
		log.Println(err)
		return v, -1, -1, -1
	}

	for rows.Next() {
		err = rows.Scan(&records)
		if err != nil {
			log.Println("[UpsVoltageGet] Scan error:")
			log.Println(err)
			return v, -1, -1, -1
		}
	}

	offset := 0
	newPg = pg
	newSz = sz
	if newSz > records {
		newPg = 1
		newSz = records
		offset = 0
	} else {
		offset = newSz * (newPg - 1)
		if offset > records {
			newPg = 1
			newSz = records
			offset = 0
		}
	}

	rows, err = db.DbConn.Query("SELECT * FROM ups ORDER BY uid DESC LIMIT " + strconv.Itoa(offset) + "," + strconv.Itoa(newSz))
	if err != nil {
		log.Println("[UpsVoltageGet] Query error:")
		log.Println(err)
		return v, -1, -1, -1
	}

	var uid uint32

	for rows.Next() {
		el := model.Volt{}
		err = rows.Scan(&uid, &el.Vcc, &el.Vin, &el.Vcin, &el.Vout, &el.Vb1, &el.Vb2, &el.Created)
		if err != nil {
			log.Println("[UpsVoltageGet] Scan error:")
			log.Println(err)
			return v, -1, -1, -1
		}
		v = append(v, el)
	}

	return v, records, newPg, newSz
}
