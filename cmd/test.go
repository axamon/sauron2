// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"net/http"
	"os"
	"time"

	"github.com/axamon/sauron2/sms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Invia sms di test",
	Long: `Per sincerarsi che le notifiche funzionino cor
	rettamente è possibile inviare un sms di test`,
	Run: func(cmd *cobra.Command, args []string) {
		cellditest := viper.GetString("Cellpertest")

		if ok := sms.Verificacellulare(cellditest); ok == false {
			fmt.Printf("cellulare nel formato errato %s\n", cellditest)
		}

		timeout := time.Duration(10 * time.Second)
		client := http.Client{
			Timeout: timeout,
		}
		resp, err := client.Get("http://google.com")
		if err != nil {
			fmt.Fprintln(os.Stdout, "errore nel contattare internet", err.Error())
			os.Exit(1)
		}

		if httpstatus := resp.StatusCode; httpstatus > 399 {
			fmt.Fprintln(os.Stdout, "errore httpstatus: ", httpstatus)
		}

		messaggio := ("Notifiche vocali correttamente funzionanti")
		result := sms.Inviasms(cellditest, messaggio)

		fmt.Println(result)
		fmt.Println("test called")
	},
}

func init() {

	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
