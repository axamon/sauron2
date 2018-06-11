package cercasid

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
)

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	err = fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
	fmt.Fprintln(os.Stderr, err.Error())
	raven.CaptureError(err, nil)
	return "", err
}

//Retrievestatus trova lo status di una call
func Retrievestatus(sid string) (status string, err error) {
	/*
		//Recupera il numero da usare con twilio
		twilionumber, err := recuperavariabile("TWILIONUMBER")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} */

	//Recupera l'accountsid di Twilio dallla variabile d'ambiente
	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")

	//Recupera il token supersegreto dalla variabile d'ambiente
	authToken, err := recuperavariabile("TWILIOAUTHTOKEN")

	url := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls/" + sid

	req, err := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	body, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	//fmt.Println(res)
	//fmt.Println(string(body))

	//Creo tipo per estrarre singolo valore da file XML
	type TwilioResponse struct {
		Status string `xml:"Call>Status"`
	}

	v := TwilioResponse{}
	err = xml.Unmarshal(body, &v)

	status = v.Status

	//fmt.Printf("Status: %s\n", v.Status)
	raven.CaptureMessage(status, nil)

	return
}
