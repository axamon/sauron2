package main

import (
	"testing"
)

func TestAddRep(t *testing.T) {

	type rep []struct {
		Nome      string
		Cognome   string
		Cellulare string
		Ok        bool
	}

	var reps rep
	reps = rep{
		{Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", Ok: true},
		{Nome: "Rep2", Cognome: "Reperibile2", Cellulare: "+391234567892", Ok: true},
		{Nome: "Rep3", Cognome: "Reperibile3", Cellulare: "+39123456789", Ok: false},
		{Nome: "Rep4", Cognome: "Reperibile4", Cellulare: "3234567893", Ok: false},
	}

	for _, Rep := range reps {

		if ok, err := addRep(Rep.Nome, Rep.Cognome, Rep.Cellulare); ok != Rep.Ok {
			t.Error("Problema a settare il Reperibile", err.Error())
		}

	}
}

func TestSetRep(t *testing.T) {
	type rep []struct {
		Giorno  string
		Cognome string
		Ok      bool
	}

	var reps rep
	reps = rep{
		{Giorno: "20180101", Cognome: "Reperibile1", Ok: true},
		{Giorno: "20180102", Cognome: "Reperibile2", Ok: true},
		{Giorno: "20180103", Cognome: "Reperibile3", Ok: false},
		{Giorno: "20180103", Cognome: "Reperibile4", Ok: false},
	}

	for _, Rep := range reps {

		if ok, err := setRep(Rep.Giorno, Rep.Cognome); ok != Rep.Ok {
			t.Error("Problema a settare il Reperibile", err.Error())
		}

	}
}

//TestIsRepSet verifica che ci sia un reperibile assegnato al giorno
func TestIsRepSet(t *testing.T) {

	type rep []struct {
		Giorno  string
		Cognome string
		Ok      bool
	}

	var reps rep
	reps = rep{
		{Giorno: "20180101", Cognome: "Reperibile1", Ok: true},
		{Giorno: "20180102", Cognome: "Reperibile2", Ok: true},
		{Giorno: "20180103", Cognome: "Reperibile3", Ok: false},
		{Giorno: "20180103", Cognome: "Reperibile4", Ok: false},
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
}

//TestIdRep verifica che gli id sul DB corrispondano
func TestIdRep(t *testing.T) {

	type rep []struct {
		Cognome string
		Ok      bool
	}

	var reps rep
	reps = rep{
		{Cognome: "Reperibile1", Ok: true},
		{Cognome: "Reperibile2", Ok: true},
		{Cognome: "Reperibile3", Ok: false},
		{Cognome: "Reperibile4", Ok: false},
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
