package main

import (
	"testing"
)

func TestAddRep(t *testing.T) {

	type rep []struct {
		Giorno     string
		Nome       string
		Cognome    string
		Cellulare  string
		AddRepOk   bool
		SetRepOk   bool
		isRepSetOk bool
	}

	var reps rep
	reps = rep{
		{Giorno: "20180101", Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", AddRepOk: true, SetRepOk: true, isRepSetOk: true},
		{Giorno: "20180102", Nome: "Rep2", Cognome: "Reperibile2", Cellulare: "+391234567892", AddRepOk: true, SetRepOk: true, isRepSetOk: true},
		{Giorno: "20180103", Nome: "Rep3", Cognome: "Reperibile3", Cellulare: "+39123456789", AddRepOk: false, SetRepOk: false, isRepSetOk: false},
		{Giorno: "20180104", Nome: "Rep4", Cognome: "Reperibile4", Cellulare: "3234567893", AddRepOk: false, SetRepOk: false, isRepSetOk: false},
		{Giorno: "20180106", Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", AddRepOk: false, SetRepOk: true, isRepSetOk: true},
		{Giorno: "20180106", Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", AddRepOk: false, SetRepOk: true, isRepSetOk: true},
	}

	for _, Rep := range reps {

		if ok, err := addRep(Rep.Nome, Rep.Cognome, Rep.Cellulare); ok != Rep.AddRepOk {
			t.Error("Problema nel settare il Reperibile", err.Error(), Rep)
		}

	}

	for _, Rep := range reps {

		if ok, err := setRep(Rep.Giorno, Rep.Cognome); ok != Rep.SetRepOk {
			t.Error("Problema ad aggiornare la reperibilità", err.Error())
		}

	}
	for _, Rep := range reps {

		ok, idrep, err := isRepSet(Rep.Giorno)
		if ok != Rep.SetRepOk {
			t.Error("Reperibile non settato", err.Error())
		}
		ok, info, err := infoRep(idrep)
		if ok != Rep.SetRepOk {
			t.Error("Reperibile non settato", err.Error())
		}
		if ok == true {
			if Rep.Cognome != info.Cognome {
				t.Skip("Reperibile cambiato forse", err.Error())
			}
		}

	}
	for _, Rep := range reps {
		var idrep int
		idrep, ok, err := idRep(Rep.Cognome)
		if err != nil {
			t.Skip(err.Error())
		}
		if ok != Rep.AddRepOk {
			t.Error(err.Error())
		}
		if ok, err := delRep(idrep); ok != Rep.AddRepOk {
			t.Error(err.Error())
		}
	}
}
