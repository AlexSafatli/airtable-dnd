package store

import (
	"../entities"
	"fmt"
	"github.com/fabioberger/airtable-go"
)

type airtableCharacter struct {
	AirtableID string             `json:"id,omitempty"`
	Fields     entities.Character `json:"fields"`
}

type airtableEncounter struct {
	AirtableID string             `json:"id,omitempty"`
	Fields     entities.Encounter `json:"fields"`
}

func OpenConnection(apiKey, baseID string) (*airtable.Client, error) {
	client, err := airtable.New(apiKey, baseID)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GetCharacters(tableName string, client *airtable.Client) []entities.Character {
	var records []airtableCharacter
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

func CreateEncounter(encounter entities.Encounter, tableName string, client *airtable.Client) (string, error) {
	if len(encounter.Name) == 0 {
		encounter.Name = fmt.Sprintf("s%d_l%d_r%d", encounter.Session,
			encounter.Level, encounter.Room)
	}
	record := airtableEncounter{Fields: encounter}
	err := client.CreateRecord(tableName, &record)
	return record.AirtableID, err
}
