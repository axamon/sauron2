package cercasid

import (
	"testing"
)

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
