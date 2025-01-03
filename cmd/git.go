package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/soarinferret/ticktask/internal/git"
)

var gitCmd = &cobra.Command{
	Use:     "git",
	Aliases: []string{},
	Short:   "Manage git repository",
	Long:    ``,
}

var gitInitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{},
	Short:   "Initialize a new git repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := git.Init()
		if err != nil {
			pExit("Error initializing git repository:", err)
		}

		pterm.Info.Println("Initialized git repository")
	},
}

var gitCloneCmd = &cobra.Command{
	Use:     "clone",
	Aliases: []string{},
	Short:   "Clone an existing ticktask git repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			pExit("Error: clone requires a single argument", nil)
		}
		err := git.Clone(args[0])
		if err != nil {
			pExit("Error cloning git repository:", err)
		}

		pterm.Info.Println("Cloned git repository")
	},
}

var gitPushCmd = &cobra.Command{
	Use:     "push",
	Aliases: []string{},
	Short:   "Push changes to the remote repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := git.Push()
		if err != nil {
			pExit("Error pushing to remote repository:", err)
		}

		pterm.Info.Println("Pushed changes to remote repository")
	},
}

var gitPullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{},
	Short:   "Pull changes from the remote repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := git.Pull()
		if err != nil {
			pExit("Error pulling from remote repository:", err)
		}

		pterm.Info.Println("Pulled changes from remote repository")
	},
}

var gitSyncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{},
	Short:   "Sync changes with the remote repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := git.Pull()
		if err != nil {
			if err.Error() != "already up-to-date" {
				pExit("Error pulling from remote repository:", err)
			}
		}

		err = git.Push()
		if err != nil {
			if err.Error() != "already up-to-date" {
				pExit("Error pushing to remote repository:", err)
			}
		}

		pterm.Info.Println("Synced changes with remote repository")
	},
}

var rootSyncCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{},
	Short:   "Sync changes with the remote repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		gitSyncCmd.Run(cmd, args)
	},
}

var gitRemoteCmd = &cobra.Command{
	Use:     "remote",
	Aliases: []string{},
	Short:   "Manage remote repositories",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		// if args, set the remote
		if len(args) == 1 {
			err := git.AddRemote("origin", args[0])
			if err != nil {
				pExit("Error setting remote:", err)
			}
		}

		// print the current remote
		remote, err := git.GetRemote()
		if err != nil {
			pExit("Error getting remote:", err)
		}

		pterm.Info.Println("Remote:", remote)
	},
}

var gitStatusCmd = &cobra.Command{
	Use:     "status",
	Aliases: []string{},
	Short:   "Show the status of the git repository",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		status, err := git.Status()
		if err != nil {
			pExit("Error getting status:", err)
		}

		pterm.Info.Println(status)
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
	rootCmd.AddCommand(rootSyncCmd)

	gitCmd.AddCommand(gitInitCmd)
	gitCmd.AddCommand(gitCloneCmd)
	gitCmd.AddCommand(gitPushCmd)
	gitCmd.AddCommand(gitPullCmd)
	gitCmd.AddCommand(gitSyncCmd)
	gitCmd.AddCommand(gitRemoteCmd)
	gitCmd.AddCommand(gitStatusCmd)
}
