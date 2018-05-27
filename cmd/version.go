package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostra la versione di sauron",
	Long:  `All software has versions. This is Sauron's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Versione:", viper.Get("Version"))
	},
}
