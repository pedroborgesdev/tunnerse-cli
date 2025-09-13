package commands

import (
	"fmt"
	"tunnerse/jobs"
	"tunnerse/logger"

	"github.com/spf13/cobra"
)

// newTunnel representa o comando "new", que cria um túnel persistente.
var infoTunnel = &cobra.Command{
	Use:                "info <tunnel_id>",
	Short:              "show tunnel information",
	DisableFlagParsing: true,
	Args:               cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		jobs.CloseKeyboardJob()
		infoRun(args[0])
	},
}

func infoRun(tunnelID string) {
	logger.Log("INFO", "trying get tunnel informations", []logger.LogDetail{}, false)

	info, err := Repo.InfoTunnel(tunnelID)
	if err != nil {
		logger.Log("FATAL", "failed to get tunnel information", []logger.LogDetail{
			{Key: "error", Value: err.Error()},
		}, false)
		return
	}

	if info == nil {
		logger.Log("INFO", "tunnel not found", nil, false)
		return
	}

	fmt.Printf(
		"\033[36mID:           \033[0m%s\n"+
			"\033[36mPort:         \033[0m%v\n"+
			"\033[36mServer:       \033[0m%s\n"+
			"\033[36mCreatedAt:    \033[0m%s\n\n"+
			"\033[32mRequests:     \033[0m%v\n"+
			"\033[38;2;255;105;180mHealthchecks: \033[0m%v\n"+
			"\033[33mWarns:        \033[0m%v\n"+
			"\033[31mErrors:       \033[0m%v\n",
		info.ID, info.Port, info.Domain, info.CreatedAt,
		info.Requests, info.Healthchecks, info.Warns, info.Errors,
	)

	// for _, t := range info {
	// 	status := "Inactive"
	// 	if t.Active {
	// 		status = "Active"
	// 	}

	// 	if !ForApp {
	// 		color := "\033[33m"
	// 		if t.Active {
	// 			color = "\033[32m"
	// 		}
	// 		fmt.Printf("%s%s\033[0m - \033[36m%s\033[0m - %s\033[0m\n", color, t.ID, t.Url, status)
	// 	} else {
	// 		fmt.Printf("id:[%s]url:[%s]status:[%s]\n", t.ID, t.Url, status)
	// 	}
	// }

}
