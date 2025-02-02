/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/soarinferret/ticktask/internal/profile"
	"github.com/soarinferret/ticktask/internal/git"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "Simple todo and task logging tool",
	Long: `A simple todo and task logging tool using the command line.
Based on the todo.txt format, this tool is designed to help
you keep track of the time you spend on tasks in as simple
and unobtrusive way as possible.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		c, _ := cmd.Flags().GetString("config")
		if c != "" {
			viper.SetConfigFile(c)
			err := viper.ReadInConfig()
			if err != nil {
				pExit("Error reading config file: ", err)
			}
		}

		p, _ := cmd.Flags().GetString("profile")
		if p != "" {
			profile.SetActiveProfile(p, false)
		}
		n, _ := cmd.Flags().GetBool("no-profile")
		if n {
			profile.UnsetActiveProfile(false)
		}

		initializeSetup()
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.commitment-clock.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// flag for the auth key

	rootCmd.PersistentFlags().StringP("config", "C", "", "Alternate configuration file to use")
	rootCmd.PersistentFlags().StringP("profile", "P", "", "Override the active profile")
	rootCmd.PersistentFlags().BoolP("no-profile", "N", false, "Ignore the active profile")

	//viper.BindPFlag("auth", rootCmd.PersistentFlags().Lookup("auth"))
	//viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
}

func globalArgsHandler(args []string) (projects []string, contexts []string) {
	// parse the arguments
	for _, arg := range args {
		if arg[0] == '+' {
			projects = append(projects, arg[1:])
		} else if arg[0] == '@' {
			contexts = append(contexts, arg[1:])
		}
	}
	return
}

func pExit(s string, err error) {
	if err != nil {
		pterm.Error.Println(s, err)
		os.Exit(1)
	}
}

func initializeSetup() {
	// offer to create the task repository or clone it
	// if the task repository does not exist
	_, err := os.Stat(viper.GetString("task_path") + "/.git")
	if os.IsNotExist(err) {
		result, _ := pterm.DefaultInteractiveConfirm.Show("Task repository does not exist. Would you like to create it?")
		if !result {
			pterm.Info.Println("Exiting...")
			os.Exit(0)
		}

		// create the task repository
		err := git.Init()
		if err != nil {
			pExit("Error creating task repository", err)
		}
		pterm.Info.Println("Initialized git repository")
	}
}
