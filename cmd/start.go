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

	"github.com/axamon/reperibili"
	"github.com/axamon/sauron/cercasid"
	"github.com/axamon/sms"
	"github.com/hpcloud/tail"
	"github.com/spf13/cobra"
)

var nagioslog, reperibilita string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//fmt.Println(nagioslog, reperibilita)
		//Verifica se il file dei log di nagios esiste
		if _, err := os.Stat(reperibilita); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, "Il file "+reperibilita+" non esiste oppure non accessibile")
			os.Exit(1)
		}

		//Inizia il tail dalla fine del file senza leggerlo dall'inizio
		var fine tail.SeekInfo
		fine.Offset = 0
		fine.Whence = 2

		//MustExist il file deve esistere Follow fa tail -f e ReOpen gestisce il logrotate
		t, err := tail.TailFile(nagioslog,
			tail.Config{
				Location:  &fine,
				MustExist: true,
				Follow:    true,
				ReOpen:    true,
			})

		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		//Per ogni nuova linea nel file
		for line := range t.Lines {

			//fmt.Println(line.Text) //per debug

			//Se la linea è di notifica la analizza se no passa oltre
			//notificabool := strings.Contains(line.Text, "NOTIFICATION")
			switch {

			case strings.Contains(line.Text, "NOTIFICATION") && strings.Contains(line.Text, "CRITICAL"):
				fmt.Println(line.Text)
				//TODO cambiare CDN con qualcosa di variabile
				reperibile, _ := reperibili.Reperibiliperpiattaforma2("CDN", reperibilita)

				TO := reperibile.Cellulare
				NOME := reperibile.Nome
				COGNOME := reperibile.Cognome
				//debug
				fmt.Println(TO, NOME, COGNOME)

				/* sid, err := reperibili.Chiamareperibile(TO, NOME, COGNOME)
				if err != nil {
					fmt.Println("Errore", err.Error())
				}
				fmt.Println(sid)
				*/

				//Cerca di chiamare il reperibile per 3 volte
				go func() {
					for n := 1; n < 5; n++ {
						//chiamo per la n volta
						fmt.Println(n) //debug
						sid, err := reperibili.Chiamareperibile(TO, NOME, COGNOME)
						if err != nil {
							fmt.Println("Errore", err.Error())
						}
						fmt.Println(sid)
						//attendi 60 secondi
						time.Sleep(90 * time.Second)
						//e verifica lo  status del sid
						status := cercasid.Retrievestatus(sid)
						//se lo status è completed esce dal loop
						if status == "completed" {
							fmt.Println("Reperibile", NOME, COGNOME, "contattattato con successo al", TO, "alle", time.Now())
							return
						}
						//se lo status è diverso da completed
						//Bisogna scalare il problema
						fmt.Println("ho provato", n, " volte e non sono riuscito a contattare il reperibile", NOME, COGNOME, TO, status)
						if n == 4 {
							for m := 1; m < 10; m++ {
								//TODO: Cambiare funzione e mettere una specifica per il servicedesk
								fmt.Println("Chiamo il numero di escalation")
								sid, err := reperibili.Chiamareperibile("NUMEEO", "UTENTE", "SERVICEDESK")
								if err != nil {
									fmt.Println("Errore", err.Error())
								}
								fmt.Println(sid)
								time.Sleep(80 * time.Second)
								status := cercasid.Retrievestatus(sid)
								if status == "completed" {
									fmt.Println("ServiceDesk contattattato con successo al alle", time.Now())
									return
								}
								fmt.Println("SD non risponde tentativo", m, time.Now())
							}
							fmt.Println("Molto grave! Neanche il SD sono riuscito a chiamare!", time.Now())
							return
						}

					}
				}()

				//esce dallo switch
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

func init() {
	rootCmd.AddCommand(startCmd)

	//variabile che punta al file log di Nagios
	//var nagioslog = flag.String("nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")
	startCmd.PersistentFlags().StringVar(&nagioslog, "nagioslog", "/var/log/nagios/nagios.log", "Nagios file di log")

	//variabile per recuperare lo storage della reperibilità
	startCmd.PersistentFlags().StringVar(&reperibilita, "reperibilita", "$GOPATH/src/github.com/axamon/sauron/sauron/reperibilita.csv", "Nagios file di log")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
