package reperibili

import (
	"fmt"
	"os"
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
		{Giorno: ieri, Nome: "Rep1", Cognome: "Reperibile1", Cellulare: "+391234567891", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: true},
		//aggiungo un secondo reperibile mai visto prima e un nuova reperibilità
		{Giorno: oggi, Nome: "Rep2", Cognome: "Reperibile2", Cellulare: "+391234567892", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: true},
		//Cellulare non buono
		{Giorno: domani, Nome: "Rep3", Cognome: "Reperibile3", Cellulare: "+39123456783", Piattaforma: "CDN", Gruppo: "Gruppo6", Ok: false},
		//Aggiungo un reperibile che in un nuova reperibilità
		{Giorno: ieri, Nome: "Rep4", Cognome: "Reperibile4", Cellulare: "+391234567894", Piattaforma: "AVS", Gruppo: "Gruppo6", Ok: true},
		//aggiungo un secondo reperibile mai visto prima e un nuova reperibilità
		{Giorno: oggi, Nome: "Rep5", Cognome: "Reperibile5", Cellulare: "+391234567895", Piattaforma: "AVS", Gruppo: "Gruppo6", Ok: true},
	}

	for _, Rep := range reps {

		err := AddRuota(Rep.Nome, Rep.Cognome, Rep.Cellulare, Rep.Piattaforma, Rep.Giorno, Rep.Gruppo)
		if err != nil {
			if Rep.Ok == false {
				t.Logf("Corretto non sia stato settato")
				break
			}
			t.Error("Problema nel settare il Reperibile", err.Error(), Rep)
		}

		_, err = IDRuota(Rep.Giorno, Rep.Piattaforma)
		if err != nil {
			t.Error("Problema nel recuperare il Reperibile", err.Error(), Rep)
		}

	}
}

func TestRecupera(t *testing.T) {
	err := os.Setenv("foo", "bar")
	if err != nil {
		t.Fatalf("Impossibile settare variabile d'ambiente. %s", err.Error())
	}

	result, err := recuperavariabile("foo")
	if err != nil {
		t.Fatalf("Impossibile recupeare variabile %s", err.Error())
	}
	if result != "bar" {
		t.Fatalf("Variabile recuperata errata %s", err.Error())
	}
	result, err = recuperavariabile("bar")
	if err != nil {
		t.Skipf("Corretto che fallisca %s", err.Error())
	}

}

func TestChiamarep(t *testing.T) {

	rep, err := GetReperibile("CDN")

	if err != nil {
		t.Errorf("Impossibile trovare reperibile: %s", err.Error())
	}

	_, err = Chiamareperibile(rep.Cellulare, rep.Nome, rep.Cognome)

	if err != nil {
		t.Errorf("Impossibile chiamare reperibile: %s", err.Error())
	}

	rep, err = GetReperibile("NULL")

	if err != nil {
		t.Skipf("Impossibile trovare reperibile: %s", err.Error())
	}

	sid, err := Chiamareperibile(rep.Cellulare, rep.Nome, rep.Cognome)
	if err != nil {
		t.Skipf("Impossibile chiamare reperibile: %s", err.Error())
	}
	if len(sid) == 0 {
		t.Failed()
	}
	fmt.Println(sid)

}
