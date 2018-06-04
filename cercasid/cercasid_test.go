package cercasid

import (
	"testing"
)

func TestCercasid(t *testing.T) {

	sid := "CA10c7de2b487b59ceb51917ab81aa2367"

	status := Retrievestatus(sid)
	if status != "busy" {
		t.Error("Errore nel recupero sid", status)
	}

}
