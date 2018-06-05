package sms

import (
	"testing"
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
			t.Error("Formato cellulare non valido", cellulare.Numcell)
		}
	}
}
