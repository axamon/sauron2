package main

import (
	"testing"
)

func TestAddRep(t *testing.T) {

	type rep []struct {
		Giorno    string
		Nome      string
		Cognome   string
		Cellulare string
		Ok        bool
	}

	var reps rep
	reps = rep{
		{Giorno: "20180101", Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", Ok: true},
		{Giorno: "20180102", Nome: "Rep2", Cognome: "Reperibile2", Cellulare: "+391234567892", Ok: true},
		{Giorno: "20180103", Nome: "Rep3", Cognome: "Reperibile3", Cellulare: "+39123456789", Ok: false},
		{Giorno: "20180104", Nome: "Rep4", Cognome: "Reperibile4", Cellulare: "3234567893", Ok: false},
	}

	for _, Rep := range reps {

		if ok, err := addRep(Rep.Nome, Rep.Cognome, Rep.Cellulare); ok != Rep.Ok {
			t.Error("Problema nel settare il Reperibile", err.Error(), Rep)
		}

	}

	for _, Rep := range reps {

		if ok, err := setRep(Rep.Giorno, Rep.Cognome); ok != Rep.Ok {
			t.Error("Problema a settare il Reperibile", err.Error())
		}

	}
	for _, Rep := range reps {

		ok, idrep, err := isRepSet(Rep.Giorno)
		if ok != Rep.Ok {
			t.Error("Reperibile non settato", err.Error())
		}
		ok, info, err := infoRep(idrep)
		if ok != Rep.Ok {
			t.Error("Reperibile non settato", err.Error())
		}
		if ok == true {
			if Rep.Cognome != info.Cognome {
				t.Error("Reperibile sbagliato", err.Error())
			}
		}

	}
	for _, Rep := range reps {
		var idrep int
		idrep, ok, err := idRep(Rep.Cognome)
		if err != nil {
			t.Skip(err.Error())
		}
		if ok != Rep.Ok {
			t.Error(err.Error())
		}
		if ok, err := delRep(idrep); ok != Rep.Ok {
			t.Error(err.Error())
		}
	}
}
