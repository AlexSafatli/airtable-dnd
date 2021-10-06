package main

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/store"
	"os"
)

var (
	commands = []string{"encounter", "items"}
)

func main() {
	// Check if arg length correct and look for command
	if len(os.Args) < 2 {
		fmt.Printf("Usage of %s:\n <command=%+v>", os.Args[0], commands)
		os.Exit(1)
	}

	// Read config file(s)
	conf := loadConfigs()

	// Check for empty credentials
	if len(conf.ApiKey) == 0 || len(conf.CampaignBase) == 0 {
		panic("No API key or campaign base ID was found to " +
			"connect to. Check your config file (this application looks" +
			" for files in ./config).")
	}

	// Open AirTable connection to campaign base
	conn, err := store.OpenConnection(conf.ApiKey, conf.CampaignBase)
	if err != nil {
		panic(err)
	}

	// Switch on command - can eventually use cobra for more advanced commands
	switch command := os.Args[1]; command {
	case "encounter":
		runEncounter(conf, conn)
	case "items":
		manageItems(conf, conn)
	default:
		// command not implemented
		fmt.Printf("Command '%s' not recognized. Needs to be one of %+v\n",
			command, commands)
	}
}
