package store

import (
	"../entities"
	"github.com/fabioberger/airtable-go"
)

type AirtableCharacter struct {
	AirtableID string
	Fields     entities.Character
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
	characters = make([]entities.Character, len(records))
	for i := range records {
		characters[i] = records[i].Fields
	}
	return characters
}
