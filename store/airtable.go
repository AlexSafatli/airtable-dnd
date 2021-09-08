package store

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/entities"
	"github.com/fabioberger/airtable-go"
)

type AirtableCharacter struct {
	AirtableID string             `json:"id,omitempty"`
	Fields     entities.Character `json:"fields"`
}

type AirtableEncounter struct {
	AirtableID string             `json:"id,omitempty"`
	Fields     entities.Encounter `json:"fields"`
}

type AirtableItem struct {
	AirtableID string        `json:"id,omitempty"`
	Fields     entities.Item `json:"fields"`
}

func OpenConnection(apiKey, baseID string) (*airtable.Client, error) {
	client, err := airtable.New(apiKey, baseID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetCharacters(tableName string, client *airtable.Client) []entities.Character {
	var records []AirtableCharacter
	if err := client.ListRecords(tableName, &records); err != nil {
		return []entities.Character{}
	}
	var characters []entities.Character
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			characters = append(characters, records[i].Fields)
		}
	}
	return characters
}

func CreateCharacter(character entities.Character, tableName string, client *airtable.Client) (string, error) {
	record := AirtableCharacter{Fields: character}
	err := client.CreateRecord(tableName, &record)
	return record.AirtableID, err
}

func UpdateCharacterByID(id string, fields map[string]interface{}, tableName string, client *airtable.Client) error {
	record := AirtableCharacter{}
	return client.UpdateRecord(tableName, id, fields, &record)
}

func GetEncounters(tableName string, client *airtable.Client) []entities.Encounter {
	var records []AirtableEncounter
	if err := client.ListRecords(tableName, &records); err != nil {
		return []entities.Encounter{}
	}
	var encounters []entities.Encounter
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			encounters = append(encounters, records[i].Fields)
		}
	}
	return encounters
}

func CreateEncounter(encounter entities.Encounter, tableName string, client *airtable.Client) (string, error) {
	if len(encounter.Name) == 0 {
		encounter.Name = fmt.Sprintf("s%d_l%d_r%d", encounter.Session,
			encounter.Level, encounter.Room)
	}
	for i, e := range GetEncounters(tableName, client) {
		if e.Name == encounter.Name {
			encounter.Name = e.Name + fmt.Sprintf("_%d", i)
		}
	}
	record := AirtableEncounter{Fields: encounter}
	err := client.CreateRecord(tableName, &record)
	return record.AirtableID, err
}

func GetItems(tableName string, client *airtable.Client) []entities.Item {
	var records = GetItemsWithIDs(tableName, client)
	var items []entities.Item
	for i := range records {
		if len(records[i].Fields.Name) > 0 {
			items = append(items, records[i].Fields)
		}
	}
	return items
}

func GetItemsWithIDs(tableName string, client *airtable.Client) []AirtableItem {
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
