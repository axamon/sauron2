package main

import (
	"flag"
	"fmt"
	"time"

	"database/sql"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

//Reperibile è la variabile con i dati personali dei reperibili
type Reperibile struct {
	gorm.Model
	Nome         string
	Cognome      string
	Cellulare    string
	Assegnazioni Assegnazione
}

//Assegnazione è la variabile con i dati relativi alla ruota di reperibilità
type Assegnazione struct {
	gorm.Model
	Piattaforma  string
	Giorno       string
	Gruppo       string
	ReperibileID uint
}

var t = time.Now()

//limite delle 7 fino alle 7 del mattino seguente il reperibile che viene visualizzato è quello del giorno prima
var limite7 = time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, t.Location())

var ieri = time.Now().Add(-24 * time.Hour).Format("20060102")
var oggi = time.Now().Format("20060102")
var domani = time.Now().Add(24 * time.Hour).Format("20060102")

var filecsv = flag.String("f", "reperibilita.csv", "Percorso del file csv per la reperibilità")
var piattaforma = flag.String("p", "CDN", "La piattaforma di cui desideri ricavare il reperibile")

var contatti []Reperibile

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

const (
	dbfile           = "reperibili.db"
	createreperibile = `
	CREATE TABLE IF NOT EXISTS reperibile (
		id	integer PRIMARY KEY AUTOINCREMENT,
		created_at	datetime,
		updated_at	datetime,
		deleted_at	datetime,
		nome	varchar ( 255 ),
		cognome	varchar ( 255 ),
		cellulare	varchar ( 255 )
	);`

	createassegnazione = `
	CREATE TABLE IF NOT EXISTS assegnazione (
		id	integer PRIMARY KEY AUTOINCREMENT,
		created_at	datetime,
		updated_at	datetime,
		deleted_at	datetime,
		piattaforma	varchar ( 255 ),
		giorno	varchar ( 255 ),
		gruppo	varchar ( 255 ),
		reperibile_id	integer
	);`
)

var db *sql.DB

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	creadb1, err := db.Prepare(createreperibile)
	checkErr(err)
	_, errcreadb1 := creadb1.Exec()
	checkErr(errcreadb1)
	creadb2, err := db.Prepare(createassegnazione)
	checkErr(err)
	_, errcreadb2 := creadb2.Exec()
	checkErr(errcreadb2)
	addreperibile, err := db.Prepare("INSERT INTO reperibile (id,nome, cognome, cellulare) VALUES (?,?, ?,?)")
	checkErr(err)
	addreperibile.Exec("1", "Alberto", "Bregliano", "+393357291533")
	addreperibile.Exec("2", "Antonio", "Gasponi", "+393357291533")
	return db
}

func main() {
	db := InitDB(dbfile)
	defer db.Close()
	id := retrieveid("Bregliano")
	fmt.Println(id)
	id = retrieveid("Gasponi")
	fmt.Println(id)

}

func retrieveid(cognome string) (id int) {
	db, err := sql.Open("sqlite3", dbfile)
	checkErr(err)
	defer db.Close()
	retrieveid, err := db.Prepare("select id from reperibile where cognome = ? limit 1")
	checkErr(err)
	row := retrieveid.QueryRow(cognome)
	err = row.Scan(&id)
	checkErr(err)
	//fmt.Println(id) //debug
	return id

}
