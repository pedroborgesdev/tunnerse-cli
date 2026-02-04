package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/dto"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/utils"

	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/jobs"
	"github.com/pedroborgesdev/tunnerse-cli/internal/cli/logger"

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

type Tunnel struct {
	ID        string
	Port      string
	Url       string
	Domain    string
	Active    bool
	CreatedAt string
}

func listRun() {
	fmt.Print(dto.Welcome)
	apiURL := "http://localhost:9988/list"
	resp, err := http.Get(apiURL)
	if err != nil {
		if utils.IsConnRefused(err) {
			logger.Log("FATAL", "Tunnerse local server is not online", []logger.LogDetail{
				{Key: "Hint", Value: "Make sure tunnerse-server is running and accessible on http://localhost:9988"},
			}, false)
		} else {
			logger.Log("FATAL", "Failed to connect to local API", []logger.LogDetail{
				{Key: "Error", Value: err.Error()},
			}, false)
		}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log("FATAL", "Failed to read server response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
		}, false)
	}

	var apiResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Tunnels []*Tunnel `json:"tunnels"`
			Count   int       `json:"count"`
		} `json:"data"`
		Status int `json:"status"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Log("FATAL", "Failed to parse server response", []logger.LogDetail{
			{Key: "Error", Value: err.Error()},
			{Key: "Response", Value: string(body)},
		}, false)
	}

	if apiResponse.Code != "success" {
		logger.Log("FATAL", "Server returned error", []logger.LogDetail{
			{Key: "Code", Value: apiResponse.Code},
			{Key: "Message", Value: apiResponse.Message},
		}, false)
	}

	tunnels := apiResponse.Data.Tunnels

	if len(tunnels) == 0 {
		logger.Log("INFO", "No tunnels found", []logger.LogDetail{}, false)
		return
	}

	activeCount := 0
	inactiveCount := 0
	for _, t := range tunnels {
		status := "Inactive"
		if t.Active {
			status = "Active"
			activeCount++
		} else {
			inactiveCount++
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

	fmt.Printf("\n\033[36mTotal tunnels:\033[0m %d | \033[32mActive:\033[0m %d | \033[33mInactive:\033[0m %d\n\n", len(tunnels), activeCount, inactiveCount)
}
