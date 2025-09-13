package commands_utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"tunnerse/config"
)

type StartOptions struct {
	LogDir  string
	LogName string
}

func StartInBackground(args []string, opts *StartOptions) (int, error) {
	cmdArgs := append([]string{"new"}, args...)

	cmd := exec.Command(os.Args[0], cmdArgs...)

	env := os.Environ()
	env = append(env,
		"TUNNERSE_BG=1",
		"TUNNERSE_TUNNEL_ID="+config.GetTunnelID(),
		"TUNNERSE_SUBDOMAIN="+strconv.FormatBool(config.GetSubdomainBool()),
		"TUNNERSE_SERVER_URL="+config.GetServerURL(),
		"TUNNERSE_ADDRESS_PORT="+strings.TrimPrefix(config.GetAddressURL(), "http://127.0.0.1:"),
	)
	cmd.Env = env

	cmd.SysProcAttr = GetSysProcAttrForBackground()

	var logPath string
	if opts != nil && opts.LogDir != "" && opts.LogName != "" {
		err := os.MkdirAll(opts.LogDir, os.ModePerm)
		if err != nil {
			return 0, fmt.Errorf("could not create log directory: %s", err.Error())
		}

		logPath = filepath.Join(opts.LogDir, opts.LogName+".log")
		logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return 0, fmt.Errorf("could not create log file: %s", err.Error())
		}
		defer logFile.Close()

		cmd.Stdout = logFile
		cmd.Stderr = logFile
	}

	if err := cmd.Start(); err != nil {
		return 0, err
	}

	return cmd.Process.Pid, nil
}
