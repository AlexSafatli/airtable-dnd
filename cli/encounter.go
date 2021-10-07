package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/fetools"
	"github.com/AlexSafatli/airtable-dnd/remote"
	"github.com/AlexSafatli/airtable-dnd/rpg"
	"github.com/fabioberger/airtable-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type encounterData struct {
	Encounter    *rpg.Encounter
	Participants []*rpg.Character
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
	var encounter encounterData
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

func createEncounter(jsonPath, monsterJsonRootPath string, monsterData []string) {
	// This code was originally found in my roleplaying-utils package and was
	// stripped out and added here as a new subcommand
	var monsters map[string]*fetools.Monster

	j, err := filepath.Glob(monsterJsonRootPath + string(os.PathSeparator) + "*.json")
	if err != nil {
		panic(err)
	}

	monsters = make(map[string]*fetools.Monster) // eventually move to API
	for _, path := range j {
		var parsed []fetools.Monster
		parsed = fetools.Get5etoolsMonsters(path)
		for i := range parsed {
			var lName = strings.ToLower(parsed[i].Name)
			if _, ok := monsters[lName]; !ok {
				monsters[lName] = &parsed[i]
			}
		}
	}

	fmt.Printf("Read %d monsters.\n", len(monsters))

	var output encounterData
	var monsterName string
	var monsterQty int
	for i, val := range monsterData {
		if i%2 == 0 {
			monsterName = val
		} else {
			monsterQty, err = strconv.Atoi(val)
			if err != nil {
				panic("Expected quantity for " + monsterName)
			}
			if m, ok := monsters[strings.ToLower(monsterName)]; ok {
				var i int
				var p rpg.Character
				p.Name = m.Name
				p.Initiative = (m.Dex - 10) / 2
				p.HP = uint(m.HP.Average)
				fmt.Printf("Monster '%s' (%d): %+v\n", m.Name, monsterQty, m)
				for i = 0; i < monsterQty; i++ {
					output.Participants = append(output.Participants, &p)
				}
			} else {
				fmt.Printf("Could not find monster '%s'\n", monsterName)
			}
		}
	}

	output.Encounter = &rpg.Encounter{}

	file, _ := json.MarshalIndent(output, "", " ")
	if err = ioutil.WriteFile(jsonPath, file, 0644); err != nil {
		panic(err)
	}

	fmt.Println("Wrote a new encounter to " + jsonPath)
}
