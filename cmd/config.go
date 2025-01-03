/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/soarinferret/ticktask/internal/config"
	"github.com/soarinferret/ticktask/internal/edit"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{},
	Short:   "View the current configuration",
	Long: `View the current configuration.
To change the default configuration, edit the file at
	` + config.GetConfigPath(true),
	Run: func(cmd *cobra.Command, args []string) {

		if cmd.Flag("edit").Changed {
			edit.OpenFile(config.GetConfigPath(false))
			return
		}

		// print raw config file
		//viper.WriteConfigTo(os.Stdout) // won't exist until viper 1.20.0 - https://github.com/spf13/viper/issues/856
		b, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
		if err != nil {
			pExit("Error writing configuration to stdout: ", err)
		}
		pterm.Println(string(b))
		return

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	configCmd.Flags().BoolP("edit", "e", false, "Open the configuration file in the default editor")

}
