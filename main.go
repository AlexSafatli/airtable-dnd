package main

import (
	"./entities"
	"./store"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type InputEncounter struct {
	Participants []*entities.Character
	Encounter    *entities.Encounter
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

func main() {
	// Check if arg length correct
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage of %s:\n <json> [submit/slots]\n", os.Args[0])
		os.Exit(1)
	}

	// Read config file(s)
	conf := loadConfigs()

	// Open AirTable connection to campaign base
	conn, err := store.OpenConnection(conf.ApiKey, conf.CampaignBase)
	if err != nil {
		panic(err)
	}

	var characters []entities.Character
	characters = store.GetCharacters(conf.TableNames.Characters, conn)
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

	// Submit to AirTable if "submit" command, otherwise do not yet
	if len(os.Args) > 2 && os.Args[2] == "submit" {
		// Save to Airtable
		var id string
		id, err = store.CreateEncounter(*encounter.Encounter, conf.TableNames.Encounters, conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Submitted '%s' to AirTable (XP %d)\n", id, encounter.Encounter.XP)
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
