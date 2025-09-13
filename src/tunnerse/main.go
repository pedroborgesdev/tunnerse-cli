package main

import (
	"flag"
	"os"
	"tunnerse/commands"
	"tunnerse/config"
	"tunnerse/utils"
)

func main() {
	defer utils.EnableInput()

	utils.DisableInput()

	forApp := flag.Bool("for-application", false, "run in app mode")
	flag.Parse()

	if *forApp {
		os.Setenv("TUNNERSE_APP", "1")
	}

	config.SetExecPath()

	commands.Execute()
}
