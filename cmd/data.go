package cmd

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/entities"
	"github.com/AlexSafatli/airtable-dnd/store"
	"github.com/fabioberger/airtable-go"
)

func getAirtableCharacters(conf configValues, conn *airtable.Client) []entities.Character {
	// Use AirTable connection to campaign base to get characters
	var characters []entities.Character
	characters = store.GetCharacters(conf.TableNames.Characters, conn)
	if len(characters) == 0 {
		errStr := fmt.Sprintf("No characters found in table '%s' on AirTable", conf.TableNames.Characters)
		panic(errStr)
	}
	return characters
}

func getAirtableItems(conf configValues, conn *airtable.Client) []entities.Item {
	// Use AirTable connection to campaign base to get party/PC items
	var items []entities.Item
	items = store.GetItems(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}

func getAirtableItemsWithIDs(conf configValues, conn *airtable.Client) []store.AirtableItem {
	// Get IDs too
	var items []store.AirtableItem
	items = store.GetItemsWithIDs(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}
