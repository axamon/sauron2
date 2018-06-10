package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/axamon/sauron2/sms"

	"database/sql"

	//serve per gestire i db sqlite
	_ "github.com/mattn/go-sqlite3"
)

//Reperibile è la variabile con i dati personali dei reperibili
type Reperibile struct {
	id           int
	Nome         string
	Cognome      string
	Cellulare    string
	Assegnazioni Assegnazione
}

//Assegnazione è la variabile con i dati relativi alla ruota di reperibilità
type Assegnazione struct {
	Piattaforma  string
	Giorno       string
	Gruppo       string
	ReperibileID uint
}

//t è il timestamp di adesso
var t = time.Now()

//limite delle 7 fino alle 7 del mattino seguente il reperibile che viene visualizzato è quello del giorno prima
var limite7 = time.Date(t.Year(), t.Month(), t.Day(), 7, 0, 0, 0, t.Location())

var ieri = time.Now().Add(-24 * time.Hour).Format("20060102")
var oggi = time.Now().Format("20060102")
var domani = time.Now().Add(24 * time.Hour).Format("20060102")

var filecsv = flag.String("f", "reperibilita.csv", "Percorso del file csv per la reperibilità")
var piattaforma = flag.String("p", "CDN", "La piattaforma di cui desideri ricavare il reperibile")

var contatti []Reperibile

const (
	//il file dove creare il DB
	dbfile = "reperibili.db"

	//Crea la tabella dei reperibili
	createreperibile = `
	CREATE TABLE IF NOT EXISTS reperibile (
		id	integer PRIMARY KEY AUTOINCREMENT,
		nome	varchar ( 255 ),
		cognome	varchar ( 255 ),
		cellulare	varchar ( 255 )
	);`

	//Crea la tabella delle reperibilità
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

//Male non fa
var db *sql.DB

//InitDB inzializza il database e restituisce la risorsa
func init() {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db nil")
	}
	creadb1, err := db.Prepare(createreperibile)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, errcreadb1 := creadb1.Exec()
	if errcreadb1 != nil {
		fmt.Println(err.Error())
	}
	creadb2, err := db.Prepare(createassegnazione)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, errcreadb2 := creadb2.Exec()
	if errcreadb2 != nil {
		fmt.Println(err.Error())
	}
	/* 	addreperibile, err := db.Prepare("INSERT INTO reperibile (id,nome, cognome, cellulare) VALUES (?,?, ?,?)")
	   	if err != nil {
	   		fmt.Println(err.Error())
	   	}
	   	_, err1 := addreperibile.Exec("1", "Alberto", "Bregliano", "+393357291533")
	   	if err1 != nil {
	   		fmt.Println(err.Error())
	   	}

	   	_, err2 := addreperibile.Exec("2", "Antonio", "Gasponi", "+393357291533")
	   	if err2 != nil {
	   		fmt.Println(err.Error())
	   	} */
}

//Opendb Opens the DB on the file system
func Opendb(pathtofile string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", pathtofile)
	if err != nil {
		//fmt.Println(err.Error())
		return nil, fmt.Errorf("Problema ad aprire il DB %s", err.Error())
	}
	return db, nil
}

func main() {
	_, err := addRep("Rep2", "Reperibile2", "+391234567892")
	if err != nil {
		fmt.Println(err)
	}
	setRep("20180609", "Reperibile2")
}

//addRep Aggiunge un reperibile al DB
func addRep(nome, cognome, cellulare string) (ok bool, err error) {
	if ok := sms.Verificacellulare(cellulare); ok != true {
		return false, fmt.Errorf("Cellulare inserito non nel formato +39(10)cifre")
	}
	db, err := Opendb(dbfile)
	defer db.Close()
	verificaprimachenonesistagia, err := db.Prepare("select count(*) from reperibile where cognome = ?")
	if err != nil {
		//fmt.Println(err.Error())
		return false, fmt.Errorf("Problema a preparare la query %s", err.Error())
	}

	var exist interface{}
	err = verificaprimachenonesistagia.QueryRow(cognome).Scan(&exist)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(exist)
	switch {
	case exist.(int64) == 0:
		addreperibile, err := db.Prepare("INSERT INTO reperibile (nome, cognome, cellulare) VALUES (?, ?,?)")
		if err != nil {
			//fmt.Println(err.Error())
			return false, fmt.Errorf("Problema a preparare la query %s", err.Error())
		}
		_, erraddrep := addreperibile.Exec(nome, cognome, cellulare)
		if erraddrep != nil {
			return false, fmt.Errorf("Impossibile inserire reperibile %s", err.Error())
		}
		return true, nil

	default:
		return false, fmt.Errorf("Impossibile inserire reperibile")
	}

}

//setRep assegna un reperibile al giorno
func setRep(giorno, cognome string) (ok bool, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()
	idrep, _, err := idRep(cognome)
	if err != nil {
		//fmt.Println(err.Error())
		return false, fmt.Errorf("Id reperibile non trovato %s", err.Error())
	}
	verifica, err := db.Prepare("select count(*) from assegnazione where giorno = ?")
	if err != nil {
		return false, fmt.Errorf("Problema a preparare la query %s", err.Error())
	}
	var exist interface{}
	err = verifica.QueryRow(giorno).Scan(&exist)
	if err != nil {
		return false, fmt.Errorf("Impossibile contare assegnazioni")
	}
	switch {
	case exist == 0:
		goto INSERISCI
	default:
		cancella, err := db.Prepare("delete from assegnazione where giorno = ?")
		if err != nil {
			//fmt.Println(err.Error())
			return false, fmt.Errorf("Problema a preparare la query %s", err.Error())
		}
		_, err = cancella.Exec(giorno)
		if err != nil {
			return false, fmt.Errorf("Impossibile cancellare")
		}
		goto INSERISCI
	}
INSERISCI:
	settaRep, err := db.Prepare("insert into assegnazione (giorno, reperibile_id) values(?,?)")
	if err != nil {
		//fmt.Println(err.Error())
		return false, fmt.Errorf("Problema a preparare la query %s", err.Error())
	}
	_, err = settaRep.Exec(giorno, idrep)
	if err != nil {
		return false, fmt.Errorf("Problema a settare il reperibile %s", err.Error())
	}
	return true, nil

}

//isRepSet informa se un Reperibile è stato impostato per il giorno e qual' è il suo id
func isRepSet(giorno string) (ok bool, reperibileID int, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()
	cercagiorno, err := db.Prepare("select reperibile_id from assegnazione where giorno = ?")

	row := cercagiorno.QueryRow(giorno)

	err = row.Scan(&reperibileID)

	switch {
	case reperibileID > 0:
		ok = true
	default:
		ok = false
	}

	return
}

//infoRep restituisce informazioni del reperibile su DB
func infoRep(idrep int) (ok bool, info Reperibile, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()
	retrieveinfo, err := db.Prepare("select nome, cognome, cellulare from reperibile where id = ? limit 1")

	row := retrieveinfo.QueryRow(idrep)

	err = row.Scan(&info.Nome, &info.Cognome, &info.Cellulare)

	switch {
	case info.Cellulare != "":
		ok = true
	default:
		ok = false
	}

	return

}

//idRep restituisce l'ID del reperibile su DB
func idRep(cognome string) (id int, ok bool, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()

	retrieveid, err := db.Prepare("select id from reperibile where cognome = ? limit 1")

	row := retrieveid.QueryRow(cognome)
	err = row.Scan(&id)

	switch {
	case id > 0:
		ok = true
	default:
		ok = false
	}

	return
}

//delRep cancella un reperibile
func delRep(idRep int) (ok bool, err error) {
	db, err := Opendb(dbfile)
	defer db.Close()
	delid, err := db.Prepare("delete from reperibile where id = ?")

	delass, err := db.Prepare("delete from assegnazione where reperibile_id = ?")

	_, err = delass.Exec(idRep)

	result, err := delid.Exec(idRep)
	affect, err := result.RowsAffected()
	switch {
	case affect == 1:
		ok = true

	default:
		ok = false
	}
	return

}
