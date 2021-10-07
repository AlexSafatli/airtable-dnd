package cmd

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

func printInitiatives(initiatives map[*entities.Character]int, directive string) {
	for _, char := range entities.RankInitiatives(initiatives).Characters() {
		// Show turn slots (factions) instead
		if directive == "slots" {
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

func runEncounter(jsonPath, directive string, conf configValues, conn *airtable.Client) {
	// Get characters
	var characters = getAirtableCharacters(conf, conn)

	// Read encounter data in as JSON; characters will be added
	var encounter InputEncounter
	f, err := os.Open(jsonPath)
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
	if directive == "submit" {
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

	printInitiatives(initiatives, directive)
}

func manageItems(jsons []string, conf configValues, conn *airtable.Client) {
	// Get party items
	var items = getAirtableItemsWithIDs(conf, conn)

	// TODO: call 5etools API microservice to get this data; for now, read JSON
	var itemMap = fetools.Get5etoolsItemMap(jsons)

	for i := range items {
		var name = strings.ToLower(items[i].Fields.Name)
		if (name == "party gold" || name == "gold") && items[i].Fields.Weight != CoinWeight {
			// Set baseline for party gold and do not go further
			var updateGoldRecord = map[string]interface{}{
				"Appr. value": 1,          // a gold coin is worth 1 gp
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
