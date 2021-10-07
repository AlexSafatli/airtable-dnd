package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/remote"
	"github.com/AlexSafatli/airtable-dnd/rpg"
	"github.com/fabioberger/airtable-go"
	"os"
)

type inputEncounter struct {
	Participants []*rpg.Character
	Encounter    *rpg.Encounter
}

func printInitiatives(initiatives map[*rpg.Character]int, directive string) {
	for _, char := range rpg.RankInitiatives(initiatives).Characters() {
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
	var encounter inputEncounter
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
		id, err = remote.CreateEncounter(*encounter.Encounter, conf.TableNames.Encounters, conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Encounter '%s' -> AirTable (XP %d, # Players %d)\n", id,
			encounter.Encounter.XP, encounter.Encounter.NumPlayers)
		return
	}

	// Roll initiative of existing characters in the encounter
	var initiatives map[*rpg.Character]int
	initiatives = make(map[*rpg.Character]int)
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
