package cli

import (
	"errors"
	"github.com/AlexSafatli/airtable-dnd/remote"
	"github.com/fabioberger/airtable-go"
	"github.com/spf13/cobra"
)

var conf configValues
var conn *airtable.Client

var cmdEncounter = &cobra.Command{
	Use:   "encounter <json> [submit/slots]",
	Short: "Manage encounters",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least a JSON path")
		}
		if len(args) >= 2 {
			if args[1] != "slots" && args[2] != "submit" {
				return errors.New("if a second argument is provided, must be a directive such as 'slots' or 'submit'")
			}
		}
		return nil
	},
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		var directive string
		if len(args) > 1 {
			directive = args[1]
		}
		runEncounter(args[0], directive, conf, conn)
	},
}

var cmdItems = &cobra.Command{
	Use:    "items <jsons>",
	Short:  "Manage items",
	Args:   cobra.MinimumNArgs(1),
	PreRun: preRun,
	Run: func(cmd *cobra.Command, args []string) {
		manageItems(args, conf, conn)
	},
}

var rootCmd = &cobra.Command{
	Use:   "airtable-dnd",
	Short: "airtable-dnd is a dungeon-mastering CLI for 5th Edition D&D campaigns",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// do nothing
	},
}

func preRun(_ *cobra.Command, _ []string) {
	// Read config file(s)
	conf = loadConfigs()

	// Check for empty credentials
	if len(conf.ApiKey) == 0 || len(conf.CampaignBase) == 0 {
		panic("No API key or campaign base ID was found to " +
			"connect to. Check your config file (this application looks" +
			" for files in ./config).")
	}

	// Open AirTable connection to campaign base
	var err error
	conn, err = remote.OpenConnection(conf.ApiKey, conf.CampaignBase)
	if err != nil {
		panic(err)
	}
}

func Execute() {
	rootCmd.AddCommand(cmdEncounter)
	rootCmd.AddCommand(cmdItems)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
