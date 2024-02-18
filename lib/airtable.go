package lib

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/rpg"
	"github.com/fabioberger/airtable-go"
)

type AirtableCharacter struct {
	AirtableID string        `json:"id,omitempty"`
	Fields     rpg.Character `json:"fields"`
}

type AirtableEncounter struct {
	AirtableID string        `json:"id,omitempty"`
	Fields     rpg.Encounter `json:"fields"`
}

type AirtableItem struct {
	AirtableID string   `json:"id,omitempty"`
	Fields     rpg.Item `json:"fields"`
}

func OpenConnection(apiKey, baseID string) (*airtable.Client, error) {
	client, err := airtable.New(apiKey, baseID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getCharacters(tableName string, client *airtable.Client) []rpg.Character {
	var records []AirtableCharacter
	if err := client.ListRecords(tableName, &records); err != nil {
		return []rpg.Character{}
	}
	var characters []rpg.Character
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			characters = append(characters, records[i].Fields)
		}
	}
	return characters
}

func CreateCharacter(character rpg.Character, tableName string, client *airtable.Client) (string, error) {
	record := AirtableCharacter{Fields: character}
	err := client.CreateRecord(tableName, &record)
	return record.AirtableID, err
}

func UpdateCharacterByID(id string, fields map[string]interface{}, tableName string, client *airtable.Client) error {
	record := AirtableCharacter{}
	return client.UpdateRecord(tableName, id, fields, &record)
}

func getEncounters(tableName string, client *airtable.Client) []rpg.Encounter {
	var records []AirtableEncounter
	if err := client.ListRecords(tableName, &records); err != nil {
		return []rpg.Encounter{}
	}
	var encounters []rpg.Encounter
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			encounters = append(encounters, records[i].Fields)
		}
	}
	return encounters
}

func createEncounter(encounter rpg.Encounter, tableName string, client *airtable.Client) (string, error) {
	if len(encounter.Name) == 0 {
		encounter.Name = fmt.Sprintf("s%d_l%d_r%d", encounter.Session,
			encounter.Level, encounter.Room)
	}
	for i, e := range getEncounters(tableName, client) {
		if e.Name == encounter.Name {
			encounter.Name = e.Name + fmt.Sprintf("_%d", i)
		}
	}
	record := AirtableEncounter{Fields: encounter}
	err := client.CreateRecord(tableName, &record)
	return record.AirtableID, err
}

func getItems(tableName string, client *airtable.Client) []rpg.Item {
	var records = getItemsWithIDs(tableName, client)
	var items []rpg.Item
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			items = append(items, records[i].Fields)
		}
	}
	return items
}

func getItemsWithIDs(tableName string, client *airtable.Client) []AirtableItem {
	var records []AirtableItem
	if err := client.ListRecords(tableName, &records); err != nil {
		panic(err)
		return []AirtableItem{}
	}
	return records
}

func UpdateItemByID(id string, fields map[string]interface{}, tableName string, client *airtable.Client) error {
	record := AirtableItem{}
	return client.UpdateRecord(tableName, id, fields, &record)
}
