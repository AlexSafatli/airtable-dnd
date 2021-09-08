package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/entities"
	"github.com/AlexSafatli/airtable-dnd/fetools"
	"github.com/AlexSafatli/airtable-dnd/store"
	"github.com/fabioberger/airtable-go"
	"os"
	"strings"
)

const CoinWeight = 1 / 50.0

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

func manageEncounter(conf configValues, conn *airtable.Client) {
	if len(os.Args) > 4 {
		fmt.Printf("Usage of %s encounter:\n encounter <json> [submit/slots]", os.Args[0])
		os.Exit(1)
	}

	// Get characters
	var characters = getAirtableCharacters(conf, conn)

	// Read encounter data in as JSON; characters will be added
	var encounter InputEncounter
	f, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	if err := json.NewDecoder(r).Decode(&encounter); err != nil {
		panic(err)
	}
	// Make adjustments to encounter object where necessary
	encounter.Encounter.NumPlayers = uint(len(characters))

	// Submit to AirTable if "submit" command, otherwise continue
	if len(os.Args) > 3 && os.Args[3] == "submit" {
		// Save to Airtable
		var id string
		id, err = store.CreateEncounter(*encounter.Encounter, conf.TableNames.Encounters, conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encounter '%s' -> AirTable (XP %d, # Players %d)\n", id,
			encounter.Encounter.XP, encounter.Encounter.NumPlayers)
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

func manageItems(conf configValues, conn *airtable.Client) {
	if len(os.Args) < 3 {
		fmt.Printf("Usage of %s items:\n items <jsons>", os.Args[0])
		os.Exit(1)
	}

	// Get party items
	var items = getAirtableItemsWithIDs(conf, conn)

	// TODO: call 5etools API microservice to get this data; for now, read JSON
	var itemMap = fetools.Get5etoolsItemMap(os.Args[2 : len(os.Args)-1])

	for i := range items {
		var name = strings.ToLower(items[i].Fields.Name)
		if name == "party gold" || name == "gold" {
			// Set baseline for party gold and do not go further
			var updateGoldRecord = map[string]interface{}{
				"Appr. value": 1,
				"Weight":      CoinWeight, // a standard coin weights 1/50 lb
			}
			if err := store.UpdateItemByID(items[i].AirtableID, updateGoldRecord, conf.TableNames.Items, conn); err != nil {
				panic(err)
			}
			fmt.Printf("Set weight baseline for party gold (name '%s') as weight %.2f lb.\n", items[i].Fields.Name, CoinWeight)
			continue
		}
		var updateRecord = make(map[string]interface{}, 2)
		if v, ok := itemMap[name]; ok {
			updateRecord["Appr. value"] = v.Value / 10 // value is in SP
			updateRecord["Weight"] = v.Weight
			if err := store.UpdateItemByID(items[i].AirtableID, updateRecord, conf.TableNames.Items, conn); err != nil {
				panic(err)
			}
			fmt.Printf("Updated item with name '%s' with value %.2f gp and weight %.2f lb.\n", v.Name, v.Value/10, v.Weight)
		}
	}
}

func main() {
	// Check if arg length correct and look for command
	if len(os.Args) < 2 {
		fmt.Printf("Usage of %s:\n <command>", os.Args[0])
		os.Exit(1)
	}

	// Read config file(s)
	conf := loadConfigs()

	// Open AirTable connection to campaign base
	conn, err := store.OpenConnection(conf.ApiKey, conf.CampaignBase)
	if err != nil {
		panic(err)
	}

	// Switch on command
	switch command := os.Args[1]; command {
	case "encounter":
		manageEncounter(conf, conn)
	case "items":
		manageItems(conf, conn)
	default:
		// command not implemented
		fmt.Printf("Command '%s' not recognized.\n", command)
	}
}
