package cercasid

import (
	"os"
	"testing"
)

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

func TestCercasid(t *testing.T) {

	sid := "CA10c7de2b487b59ceb51917ab81aa2367"

	status, err := Retrievestatus(sid)
	if err != nil {
		t.Errorf("Errore %s", err.Error())
	}
	if status != "busy" {
		t.Skip("Errore nel recupero sid", status)
	}

}
