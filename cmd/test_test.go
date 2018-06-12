package cmd

import (
	"testing"

	"github.com/axamon/sauron2/sms"
)

func TestTest(t *testing.T) {
	type formati []struct {
		Numcell string
		Valido  bool
	}

	var numeri formati
	numeri = formati{
		//manca il +39
		{Numcell: "3353458144", Valido: false},
		//non italiano
		{Numcell: "+383353458144", Valido: false},
		//non italiano e mancano cifre
		{Numcell: "+38335345814", Valido: false},
		//Numero valido
		{Numcell: "+393353458144", Valido: true},
	}
	for _, cellulare := range numeri {
		if ok := sms.Verificacellulare(cellulare.Numcell); ok != cellulare.Valido {
			t.Error("Formato cellulare non valido", cellulare.Numcell)
		}
	}

}
