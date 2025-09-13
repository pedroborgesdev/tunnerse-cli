package commands

import (
	"fmt"
	"tunnerse/jobs"
	"tunnerse/logger"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var listTunnel = &cobra.Command{
	Use:                "list",
	Short:              "list all tunnels",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		listRun()
	},
}

func listRun() {
	logger.Log("INFO", "trying get tunnels informations", []logger.LogDetail{}, false)

	tunnels, err := Repo.ListTunnels()
	if err != nil {
		logger.Log("FATAL", "failed to get tunnels", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
	}

	if len(tunnels) == 0 {
		logger.Log("INFO", "no tunels found", []logger.LogDetail{}, false)
	}

	for _, t := range tunnels {
		status := "Inactive"
		if t.Active {
			status = "Active"
		}

		if !ForApp {
			color := "\033[33m"
			if t.Active {
				color = "\033[32m"
			}
			fmt.Printf("%s%s\033[0m - \033[36m%s\033[0m - %s\033[0m\n", color, t.ID, t.Url, status)
		} else {
			fmt.Printf("id:[%s]url:[%s]status:[%s]\n", t.ID, t.Url, status)
		}
	}

}
