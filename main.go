package main

import (
	"./entities"
	"./store"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/evalphobia/go-config-loader"
	"os"
)

const (
	confType = "toml"
	basePath = "config"
)

type InputEncounter struct {
	Participants []*entities.Character
	Encounter    *entities.Encounter
}

func main() {
	// Check if arg length correct
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage of %s:\n <json> [submit/slots]\n", os.Args[0])
		os.Exit(1)
	}

	// Read config file(s)
	conf := config.NewConfig()
	if err := conf.LoadConfigs(basePath, confType); err != nil {
		panic(err)
	}

	// Open AirTable connection to campaign base
	conn, err := store.OpenConnection(conf.ValueString("database.api_key"), conf.ValueString("campaign.base_id"))
	if err != nil {
		panic(err)
	}

	var characters []entities.Character
	characters = store.GetCharacters(conf.ValueString("characters.table_name"), conn)
	if len(characters) == 0 {
		panic("No characters found on AirTable")
	}

	// Read encounter data in as JSON; characters will be added to
	var encounter InputEncounter
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	if err := json.NewDecoder(r).Decode(&encounter); err != nil {
		panic(err)
	}

	// Submit to AirTable if "submit" command, otherwise get additional info
	if len(os.Args) > 2 && os.Args[2] == "submit" {
		// Save to Airtable
		_, err = store.CreateEncounter(*encounter.Encounter,
			conf.ValueString("encounters.table_name"), conn)
		if err != nil {
			panic(err)
		}
		return
	}

	// Roll initiative of existing characters in the encounter
	var initiatives map[*entities.Character]int
	initiatives = make(map[*entities.Character]int)
	for _, char := range encounter.Participants {
		initiatives[char] = RollDice(1, 20) + char.Initiative
	}

	// Add PCs, get their initiatives
	for i := range characters {
		var initiative int
		characters[i].Affiliated = true
		encounter.Participants = append(encounter.Participants, &characters[i])
		fmt.Printf("Enter Initiative for %s: ", characters[i].Name)
		if _, err := fmt.Scanf("%d", &initiative); err != nil {
			initiative = RollDice(1, 20) + characters[i].Initiative
		}
		initiatives[&characters[i]] = initiative
	}
	printInitiatives(initiatives)
}

func printInitiatives(initiatives map[*entities.Character]int) {
	for _, char := range entities.RankInitiatives(initiatives).Characters() {
		// Show turn slots (factions) instead
		if len(os.Args) > 2 && os.Args[2] == "slots" {
			var name string
			if char.Affiliated {
				name = char.Name
			} else {
				name = "Enemy"
			}
			fmt.Printf("%-8s\t%d\n", name, initiatives[char])
		} else {
			fmt.Printf("%-8s\t%d\t%d\n", char.Name, initiatives[char], char.HP)
		}
	}
}
