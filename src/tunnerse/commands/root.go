package commands

import (
	"tunnerse/config"
	"tunnerse/database"
	"tunnerse/servers"

	"github.com/spf13/cobra"
)

var (
	db     = database.NewDatabase()
	Repo   = database.NewActionsRepository(db)
	Server = servers.NewServerService(Repo)
	ForApp bool
)

var rootCmd = &cobra.Command{
	Use:   "tunnerse",
	Short: "Expose your on-premises application to the entire internet via reverse tunnels",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.PersistentFlags().BoolVar(&ForApp, "for-application", false, "flag to run program with Tunerse desktop portability\nwarn: must be used before arguments\n")
	// rootCmd.PersistentFlags().MarkHidden("for-application")

	if config.GetEnvApplicationRunning() {
		ForApp = true
	}

	rootCmd.AddCommand(quickTunnel)
	rootCmd.AddCommand(newTunnel)
	rootCmd.AddCommand(logsTunnel)
	rootCmd.AddCommand(killTunnel)
	rootCmd.AddCommand(listTunnel)
	rootCmd.AddCommand(infoTunnel)
	rootCmd.Execute()
}
