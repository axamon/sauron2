package reperibili

import (
	"testing"
)

func TestAddrep(t *testing.T) {

	type rep []struct {
		ID          int
		Nome        string
		Cognome     string
		Cellulare   string
		Piattaforma string
		Giorno      string
		Gruppo      string
		Ok          bool
	}

	var reps rep
	reps = rep{
		//Aggiungo un reperibile che in un nuova reperibilità
		{Giorno: "20180101", Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: true},
		//aggiungo un secondo reperibile mai visto prima e un nuova reperibilità
		{Giorno: "20180102", Nome: "Rep2", Cognome: "Reperibile2", Cellulare: "+391234567892", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: true},
		//Cellulare non buono
		//{Giorno: "20180102", Nome: "Rep3", Cognome: "Reperibile3", Cellulare: "+39123456783", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: false},
	}

	for _, Rep := range reps {

		err := AddRuota(Rep.Nome, Rep.Cognome, Rep.Cellulare, Rep.Piattaforma, Rep.Giorno, Rep.Gruppo)
		if err != nil {
			t.Error("Problema nel settare il Reperibile", err.Error(), Rep)
		}

		_, err = IDRuota(Rep.Giorno, Rep.Piattaforma)
		if err != nil {
			t.Error("Problema nel recuperare il Reperibile", err.Error(), Rep)
		}

	}
}
