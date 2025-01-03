/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/soarinferret/ticktask/internal/profile"
)

var profileCmd = &cobra.Command{
	Use:     "profile",
	Aliases: []string{"p"},
	Short:   "Manage local profiles",
	Long:    ``,
}

var profileSwitchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"s", "sw"},
	Short:   "Switch to a different context",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := profile.SetActiveProfile(args[0], true)
		if err != nil {
			pExit("Failed to switch profile:", err)
		}
		pterm.Info.Println("Switched to profile: ", args[0])
	},
}

var profileListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all profiles",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		printProfileTable(profile.GetProfiles())
	},
}

var profileAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a", "create"},
	Short:   "Add a new profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		projects, _ := cmd.Flags().GetStringSlice("projects")
		contexts, _ := cmd.Flags().GetStringSlice("contexts")
		active, _ := cmd.Flags().GetBool("active")
		name, _ := cmd.Flags().GetString("name")

		p := profile.AddProfile(name, active, projects, contexts)

		printProfileTable([]profile.Profile{*p})
	},
}

var profileEditCmd = &cobra.Command{
	Use:     "edit",
	Aliases: []string{"e", "modify"},
	Short:   "Edit an existing profile",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {

		projects, contexts := globalArgsHandler(args)

		name, _ := cmd.Flags().GetString("name")
		projectsFlag, _ := cmd.Flags().GetStringSlice("projects")
		contextsFlag, _ := cmd.Flags().GetStringSlice("contexts")

		// merge projects and contexts
		projects = append(projects, projectsFlag...)
		contexts = append(contexts, contextsFlag...)

		p, err := profile.EditProfile(name, projects, contexts)
		if err != nil {
			pExit("Failed to edit profile:", err)
		}

		printProfileTable([]profile.Profile{*p})
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)

	profileCmd.AddCommand(profileSwitchCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	profileAddCmd.Flags().StringSliceP("projects", "p", []string{}, "Projects to add to the profile")
	profileAddCmd.Flags().StringSliceP("contexts", "c", []string{}, "Contexts to add to the profile")
	profileAddCmd.Flags().BoolP("active", "a", false, "Set this profile as the active profile")
	profileAddCmd.Flags().StringP("name", "n", "", "The name of the profile to add")
	profileAddCmd.MarkFlagRequired("name")

	profileEditCmd.Flags().StringSliceP("projects", "p", []string{}, "Projects to add to the profile")
	profileEditCmd.Flags().StringSliceP("contexts", "c", []string{}, "Contexts to add to the profile")
	profileEditCmd.Flags().StringP("name", "n", "", "Name of profile to edit")
	profileEditCmd.MarkFlagRequired("name")

}

func printProfileTable(profiles []profile.Profile) {
	// print profiles in a table
	profileData := [][]string{}

	// add header
	profileData = append(profileData, []string{"Name", "Active", "Projects", "Contexts"})

	// add profile data
	for _, p := range profiles {
		profileData = append(
			profileData,
			[]string{
				p.Name,
				strconv.FormatBool(p.Active),
				strings.Join(p.Projects, ", "),
				strings.Join(p.Contexts, ", "),
			},
		)
	}
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(profileData).Render()
}
