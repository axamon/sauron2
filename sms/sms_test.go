package sms

import (
	"fmt"
	"os"
	"testing"

	"github.com/getsentry/raven-go"
)

func TestVerificacellulare(t *testing.T) {
	type formati []struct {
		Numcell string
		Valido  bool
	}

	var numeri formati
	numeri = formati{
		{Numcell: "3353458144", Valido: false},
		{Numcell: "+383353458144", Valido: false},
		{Numcell: "+38335345814", Valido: false},
		{Numcell: "+393353458144", Valido: true},
	}
	for _, cellulare := range numeri {
		if ok := Verificacellulare(cellulare.Numcell); ok != cellulare.Valido {
			err := fmt.Errorf("Formato cellulare non valido", cellulare.Numcell)
			raven.CaptureErrorAndWait(err)
			t.Error(err)
		}
	}
}

func TestRecuperaVariabile(t *testing.T) {
	type variabile []struct {
		Nomevar string
		Value   string
	}
	if err := os.Setenv("zoo", "balu"); err != nil {
		t.Error(err.Error())
	}

	var element variabile
	element = variabile{
		{Nomevar: "zoo", Value: "balu"},
	}
	for _, Ele := range element {
		if result, err := recuperavariabile(Ele.Nomevar); result != Ele.Value {
			t.Error(err.Error())
		}
	}

}

func TestInviasms(t *testing.T) {
	type formati []struct {
		Numcell string
		Valido  bool
	}

	var numeri formati
	numeri = formati{
		{Numcell: "3353458144", Valido: false},
		{Numcell: "+383353458144", Valido: false},
		{Numcell: "+38335345814", Valido: false},
		{Numcell: "+393353458144", Valido: true},
	}
	for _, num := range numeri {
		if result := Inviasms(num.Numcell, "prova"); result != "201 CREATED" {
			t.Skip("Bisogna settare le variabili")
		}
	}
}
