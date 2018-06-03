package cmd

import (
	"fmt"
	"testing"
	"time"
)

func TestPrimo(t *testing.T) {
	type fobtest []struct {
		Date         time.Time
		OraInizioFob int
		Fob          bool
	}
	var tests fobtest

	tests = fobtest{
		{Date: time.Date(2018, 6, 4, 6, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
		{Date: time.Date(2018, 6, 4, 6, 59, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
		{Date: time.Date(2018, 6, 4, 7, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: false},
		{Date: time.Date(2018, 6, 4, 18, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
		{Date: time.Date(2018, 6, 3, 7, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
		{Date: time.Date(2018, 6, 2, 7, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
		{Date: time.Date(2018, 6, 4, 12, 0, 0, 0, time.Local), OraInizioFob: 18, Fob: false},
		{Date: time.Date(2018, 6, 4, 18, 30, 0, 0, time.Local), OraInizioFob: 18, Fob: true},
	}

	for _, test := range tests {
		if ok := isfob(test.Date, test.OraInizioFob); ok != test.Fob {
			t.Error("Test failed", ok)
		}
	}
}

func Exampleisfob() {
	date := time.Date(2018, 6, 4, 6, 0, 0, 0, time.Local)
	ok := isfob(date, 18)
	fmt.Println(date.Weekday())
	fmt.Println(date.Hour())
	fmt.Println(ok)
	//Output:
	//Monday
	//6
	//true
}
