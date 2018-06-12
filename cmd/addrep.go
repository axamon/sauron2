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
	"log"
	"os"

	"github.com/axamon/sauron2/reperibili"
	"github.com/spf13/cobra"
)

// addrepCmd represents the addrep command
var addrepCmd = &cobra.Command{
	Use:   "addrep",
	Short: "Inserisce le reperibilità",
	Long: `Sintassi:

sauron2 addrep <nome> <cognome> <cellulare> <piattaforma> <giorno> <gruppo>`,
	Run: func(cmd *cobra.Command, args []string) {
		nome := args[0]
		cognome := args[1]
		cellulare := args[2]
		piattaforma := args[3]
		giorno := args[4]
		gruppo := args[5]
		fmt.Println(nome, cognome, cellulare, piattaforma, giorno, gruppo)
		err := reperibili.AddRuota(nome, cognome, cellulare, piattaforma, giorno, gruppo)
		if err != nil {
			err = fmt.Errorf("Impossibile aggiungere reperibilità %s", err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
			log.Fatal(err.Error())
		}
		rep, err := reperibili.GetReperibile("CDN")
		if err != nil {
			err = fmt.Errorf("Impossibile trovare reperibile odierno: %s\n Lanciare addrep per aggngere", err.Error())
			fmt.Fprintln(os.Stderr, err.Error())
			log.Fatal(err.Error())
		}
		fmt.Println(rep.Cognome)
		fmt.Println("Aggiunta reperibilità terminata")
	},
}

func init() {
	rootCmd.AddCommand(addrepCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addrepCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addrepCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
