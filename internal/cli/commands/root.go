package commands

import (
	"fmt"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/dto"

	"github.com/spf13/cobra"
)

var (
	ForApp bool
)

var rootCmd = &cobra.Command{
	Use:   "tunnerse",
	Short: "Expose your on-premises application to the entire internet via reverse tunnels",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(dto.Welcome)
		fmt.Print(dto.Commands)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	rootCmd.AddCommand(quickTunnel)
	rootCmd.AddCommand(newTunnel)
	rootCmd.AddCommand(logsTunnel)
	rootCmd.AddCommand(killTunnel)
	rootCmd.AddCommand(delTunnel)
	rootCmd.AddCommand(listTunnel)
	rootCmd.AddCommand(infoTunnel)
	rootCmd.Execute()
}
