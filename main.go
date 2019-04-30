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
	if len(os.Args) != 3 {
		fmt.Printf("Usage of %s:\n <submit/slots/init> <json>\n", os.Args[0])
		os.Exit(1)
	}

	conf := config.NewConfig()
	if err := conf.LoadConfigs(basePath, confType); err != nil {
		panic(err)
	}
	conn, err := store.OpenConnection(conf.ValueString("database.api_key"), conf.ValueString("campaign.base_id"))
	if err != nil {
		panic(err)
	}

	var characters []entities.Character
	characters = store.GetCharacters(conf.ValueString("characters.table_name"), conn)

	if len(characters) == 0 {
		panic("No characters found")
	}

	var initiatives map[*entities.Character]int
	initiatives = make(map[*entities.Character]int)

	// Read encounter data in as JSON; characters will be added to
	var encounter InputEncounter
	f, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	if err := json.NewDecoder(r).Decode(&encounter); err != nil {
		panic(err)
	}

	// Submit to AirTable if "submit" command, otherwise get additional info
	if os.Args[1] == "submit" {
		// Save to Airtable
		_, err = store.CreateEncounter(*encounter.Encounter,
			conf.ValueString("encounters.table_name"), conn)
		if err != nil {
			panic(err)
		}
		return
	}

	// Roll initiative of existing characters in the encounter
	for _, char := range encounter.Participants {
		initiative := RollDice(1, 20) + char.Initiative
		initiatives[char] = initiative
	}

	// Add PCs, get their initiatives
	for i := range characters {
		var initiative int
		characters[i].Affiliated = true
		encounter.Participants = append(encounter.Participants, &characters[i])
		fmt.Printf("Enter Initiative for %s: ", characters[i].Name)
		if _, err := fmt.Scanf("%d", &initiative); err != nil {
			panic(err)
		}
		initiatives[&characters[i]] = initiative
	}

	for _, char := range entities.RankInitiatives(initiatives).Characters() {
		// Show turn slots (factions) instead
		if os.Args[1] == "slots" {
			var name string
			if char.Affiliated {
				name = char.Name
			} else {
				name = "Enemy"
			}
			fmt.Printf("%-8s\t%d\t%t\n", name, initiatives[char], char.Affiliated)
		} else {
			fmt.Printf("%-8s\t%d\t%d\n", char.Name, initiatives[char], char.HP)
		}
	}
}
