package cli

import (
	"errors"
	"github.com/AlexSafatli/airtable-dnd/config"
	"github.com/AlexSafatli/airtable-dnd/lib"
	"github.com/fabioberger/airtable-go"
	"github.com/spf13/cobra"
)

var conf config.Values
var conn *airtable.Client

var rootEncounterCmd = &cobra.Command{
	Use:   "encounter <command>",
	Short: "Manage encounters",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// do nothing
	},
}

var cmdEncounterCreate = &cobra.Command{
	Use:   "create <output_json> <input_5etools_json_folder> [monsterName] [monsterQty] ...",
	Short: "Create encounters",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an output JSON path")
		} else if len(args) < 2 {
			return errors.New("requires an input JSON file folder with 5etools monster data")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		lib.CreateEncounter(args[0], args[1], args[2:])
	},
}

var cmdEncounterRun = &cobra.Command{
	Use:   "run <json> [submit/slots]",
	Short: "Run an encounter",
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
		lib.RunEncounter(args[0], directive, conf, conn)
	},
}

var cmdItems = &cobra.Command{
	Use:    "items <5etools_item_jsons>",
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
	conf = config.LoadConfigs()

	// Check for empty credentials
	if len(conf.ApiKey) == 0 || len(conf.CampaignBase) == 0 {
		panic("No API key or campaign base ID was found to " +
			"connect to. Check your config file (this application looks" +
			" for files in ./config).")
	}

	// Open AirTable connection to campaign base
	var err error
	conn, err = lib.OpenConnection(conf.ApiKey, conf.CampaignBase)
	if err != nil {
		panic(err)
	}
}

func Execute() {
	rootEncounterCmd.AddCommand(cmdEncounterRun)
	rootEncounterCmd.AddCommand(cmdEncounterCreate)
	rootCmd.AddCommand(rootEncounterCmd)
	rootCmd.AddCommand(cmdItems)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
