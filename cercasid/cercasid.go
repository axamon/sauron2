package cercasid

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func recuperavariabile(variabile string) (result string, err error) {
	if result, ok := os.LookupEnv(variabile); ok && len(result) != 0 {
		return result, nil
	}
	return "", fmt.Errorf("la variabile %s non esiste o è vuota", variabile)
}

//Retrievestatus trova lo status di una call
func Retrievestatus(sid string) (status string) {
	/*
		//Recupera il numero da usare con twilio
		twilionumber, err := recuperavariabile("TWILIONUMBER")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} */

	//Recupera l'accountsid di Twilio dallla variabile d'ambiente
	accountSid, err := recuperavariabile("TWILIOACCOUNTSID")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	//Recupera il token supersegreto dalla variabile d'ambiente
	authToken, err := recuperavariabile("TWILIOAUTHTOKEN")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	url := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Calls/" + sid

	req, _ := http.NewRequest("GET", url, nil)

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, errres := http.DefaultClient.Do(req)
	if errres != nil {
		log.Fatal(errres)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))

	//Creo tipo per estrarre singolo valore da file XML
	type TwilioResponse struct {
		Status string `xml:"Call>Status"`
	}

	v := TwilioResponse{}
	errstat := xml.Unmarshal(body, &v)
	if errstat != nil {
		fmt.Printf("error: %v", err)
		return
	}

	//fmt.Printf("Status: %s\n", v.Status)

	return v.Status
}
