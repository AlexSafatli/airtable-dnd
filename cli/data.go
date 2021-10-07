package cli

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/remote"
	"github.com/AlexSafatli/airtable-dnd/rpg"
	"github.com/fabioberger/airtable-go"
)

func getAirtableCharacters(conf configValues, conn *airtable.Client) []rpg.Character {
	// Use AirTable connection to campaign base to get characters
	var characters []rpg.Character
	characters = remote.GetCharacters(conf.TableNames.Characters, conn)
	if len(characters) == 0 {
		errStr := fmt.Sprintf("No characters found in table '%s' on AirTable", conf.TableNames.Characters)
		panic(errStr)
	}
	return characters
}

func getAirtableItems(conf configValues, conn *airtable.Client) []rpg.Item {
	// Use AirTable connection to campaign base to get party/PC items
	var items []rpg.Item
	items = remote.GetItems(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}

func getAirtableItemsWithIDs(conf configValues, conn *airtable.Client) []remote.AirtableItem {
	// Get IDs too
	var items []remote.AirtableItem
	items = remote.GetItemsWithIDs(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}
