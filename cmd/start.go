// Copyright © 2018 Alberto Bregliano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/axamon/sauron2/cercasid"
	"github.com/axamon/sauron2/reperibili"
	"github.com/axamon/sauron2/sms"
	"github.com/hpcloud/tail"
	"github.com/spf13/cobra"
)

//Crea variabile per assegnare i valori presi dal file di configurazione viper
var nagioslog, reperibilita, nagiosuser string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Avvia le notifiche a voce per il reperibile in turno",
	Long: `Sauron2 ascolta in tail sul file di nagios per pattern specifici
	Se li riscontra allora contatta il reperibile in turno.`,
	Run: func(cmd *cobra.Command, args []string) {

		//Recupera ora inzio fob dal file di congigurazione
		foborainizio := viper.GetInt("foborainizio")
		fmt.Printf("Ora inzio FOB: %d\n", foborainizio)

		if fob := isfob(time.Now(), foborainizio); fob == true {
			fmt.Println("Siamo in FOB. Notifiche vocali attive!")
		}

		//Recupera file reperibilita.csv dal file di congigurazione
		reperibilita := viper.GetString("Reperibilita")

		fmt.Println(reperibilita) //Debug

		//Recupera piattaforma dal file di congigurazione
		piattaforma := viper.GetString("piattaforma")

		fmt.Println(piattaforma) //Debug

		//Verifica se il file della reperibilita esiste e se è raggiungibile
		if _, err := os.Stat(reperibilita); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Il file "+reperibilita+" non esiste oppure non accessibile")
			os.Exit(1)
		}

		//Verfica esistenza reperibile per oggi e domani
		//TODO: creare sistema gestione reperibili serio

		//Recupera file dei log nagios dal file di congigurazione
		nagioslog := viper.GetString("Nagioslogfile")

		fmt.Println(nagioslog) //Debug

		//Recupera il nome dell'utente nagios di servizio per le notifiche
		nagiosuser := viper.GetString("Nagiosuser")

		//Verifica se il file deli log nagios esiste e se è raggiungibile
		if _, err := os.Stat(nagioslog); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Il file "+nagioslog+" non esiste oppure non accessibile")
			os.Exit(1)
		}

		//Inizia il tail dalla fine del file leggendolo dalla fine
		var fine tail.SeekInfo
		fine.Offset = 0
		fine.Whence = 2

		//MustExist il file deve esistere Follow fa tail -f e ReOpen gestisce il logrotate
		//il file nagioslog
		t, err := tail.TailFile(nagioslog,
			tail.Config{
				Location:  &fine,
				MustExist: true,
				Follow:    true,
				ReOpen:    true,
			})

		//Se ci sono problemi ad accedere al file nagioslog esce.
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			//TODO: verificare che sia opportuno questo exit, fa spegnere sauron...
			os.Exit(1)
		}

		//Per ogni nuova linea nel file
	LINE:
		for line := range t.Lines {

			//fmt.Println(line.Text) //per debug

			//Se la linea non contiene il nome dello user di nagios ricomincia da LINE
			if !strings.Contains(line.Text, nagiosuser) {
				continue LINE
			}

			switch {

			//Se la linea è di notifica, la analizza, se no passa oltre
			case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "CRITICAL"):
				fmt.Println(line.Text) //debug

				//recupera il reperibile odierno tenendo conta degli orari di reperibilità
				reperibile, err := reperibili.Reperibiliperpiattaforma2(viper.GetString("piattaforma"), reperibilita)
				if err != nil {
					fmt.Fprintln(os.Stdout, err.Error())
				}

				TO := reperibile.Cellulare
				NOME := reperibile.Nome
				COGNOME := reperibile.Cognome

				fmt.Println(TO, NOME, COGNOME) //debug

				//Se non siamo in FOB non fare nulla
				//L'ora di demarcazione del fob è impostabile nel file di configurazione

				if ok := isfob(time.Now(), foborainizio); ok == false {
					fmt.Println("Siamo in orario base quindi niente notifiche")
					continue LINE
				}

				//Cerca di chiamare il reperibile per tot volte
				//TODO: se il problema rientra smettere di chiamare
				go func() {
					for n := 1; n < 5; n++ {
						//chiamo per la n volta
						fmt.Println(n) //debug
						sid, err := reperibili.Chiamareperibile(TO, NOME, COGNOME)
						if err != nil {
							fmt.Println("Errore", err.Error())
						}
						fmt.Println(sid)
						//attendi 90 secondi
						time.Sleep(90 * time.Second)
						//e verifica lo  status del sid
						status := cercasid.Retrievestatus(sid)
						//se lo status è completed esce dalla gooutine
						if status == "completed" {
							fmt.Println(time.Now().Format(time.RFC3339), "Reperibile", NOME, COGNOME, "contattattato con successo al", TO)
							return
						}
						//se lo status è diverso da completed
						//Bisogna scalare il problema
						fmt.Println(time.Now().Format(time.RFC3339), "ho provato", n, " volte e non sono riuscito a contattare il reperibile", NOME, COGNOME, TO, status)
						if n == 4 {
							for m := 1; m < 10; m++ {
								//TODO: Cambiare funzione e mettere una specifica per il servicedesk
								fmt.Println(time.Now().Format(time.RFC3339), "Chiamo il numero di escalation")
								sid, err := reperibili.Chiamareperibile(viper.GetString("numservicedesk"), "UTENTE", "SERVICEDESK")
								if err != nil {
									fmt.Println("Errore", err.Error())
								}
								fmt.Println(sid)
								time.Sleep(80 * time.Second)
								status := cercasid.Retrievestatus(sid)
								//se lo status è completed esce dalla gooutine
								if status == "completed" {
									fmt.Println(time.Now().Format(time.RFC3339), "ServiceDesk contattattato con successo")
									return
								}
								fmt.Println(time.Now().Format(time.RFC3339), "SD non risponde tentativo", m)
							}
							//Esce dalla goroutine senza essere riuscito a chiamare il Servicedesk
							//TODO: Verificare le politiche di escalation
							fmt.Println(time.Now().Format(time.RFC3339), "Molto grave! Neanche il SD sono riuscito a chiamare!")
							return
						}

					}
				}()

				//esce dallo switch e permette così di gestire nuove notifiche
				//il che potrebbe essere un problema se arrivano molteplici notifiche per piattaforma
				//Forse sarebbe oppotuno limitare il numero di chiamate a 1 per volta
				break

			case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "OK"):
				//Se ok allora manda solo sms senza chiamata
				fmt.Println("ricevuto OK") //per debug

				reperibile, err := reperibili.Reperibiliperpiattaforma2("CDN", reperibilita)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
				}

				TO := reperibile.Cellulare
				pezzi := strings.Split(line.Text, ";")
				messaggio := "Su " + pezzi[1] + " servizio " + pezzi[2] + " " + pezzi[3]

				go sms.Inviasms(TO, messaggio)
				//esce dallo switch
				break

			default:
				//fmt.Println("debug")
				//esce dallo switch
				break
			}
		}

		fmt.Println("start called")
	},
}

func isfob(ora time.Time, foborainizio int) (ok bool) {
	//ora := time.Now()
	giorno := ora.Weekday()
	//Partiamo che non siamo in FOB
	ok = false

	switch giorno {
	//Se è sabato siamo in fob
	case time.Saturday:
		//fmt.Println("E' sabato")
		ok = true
	//Se è domenica siamo in fob
	case time.Sunday:
		//fmt.Println("E' Domenica")
		ok = true
	//Se invece è un giorno feriale dobbiamo vedere l'orario
	default:
		//se è dopo le 18 siamo in fob
		//Si avviso il reperibile mezz'ora prima se è un problema si può cambiare
		//Recupero l'ora del FOB dal file di configurazione
		if ora.Hour() >= foborainizio {
			//fmt.Println("Giorno feriale", viper.GetInt("foborainizio"))
			ok = true
			return ok
		}
		//se è prima delle 7 allora siamo in fob
		if ora.Hour() < 7 {
			ok = true
		}
	}
	//Ritorna ok che sarà true o false a seconda se siamo in FOB o no
	return ok
}

func init() {
	rootCmd.AddCommand(startCmd)

	//variabile che punta al file log di Nagios
	//var nagioslog = flag.String("nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")
	//startCmd.PersistentFlags().StringVar(&nagioslog, "nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")

	//variabile per recuperare lo storage della reperibilità
	//startCmd.PersistentFlags().StringVar(&reperibilita, "reperibilita", "$GOPATH/src/github.com/axamon/sauron/sauron/reperibilita.csv", "Nagios file di log")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
