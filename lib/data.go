package lib

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/config"
	"github.com/AlexSafatli/airtable-dnd/rpg"
	"github.com/fabioberger/airtable-go"
)

func GetAirtableCharacters(conf config.Values, conn *airtable.Client) []rpg.Character {
	// Use AirTable connection to campaign base to get characters
	var characters []rpg.Character
	characters = getCharacters(conf.TableNames.Characters, conn)
	if len(characters) == 0 {
		errStr := fmt.Sprintf("No characters found in table '%s' on AirTable", conf.TableNames.Characters)
		panic(errStr)
	}
	return characters
}

func GetAirtableItems(conf config.Values, conn *airtable.Client) []rpg.Item {
	// Use AirTable connection to campaign base to get party/PC items
	var items []rpg.Item
	items = getItems(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}

func GetAirtableItemsWithIDs(conf config.Values, conn *airtable.Client) []AirtableItem {
	// Get IDs too
	var items []AirtableItem
	items = getItemsWithIDs(conf.TableNames.Items, conn)
	if len(items) == 0 {
		errStr := fmt.Sprintf("No items found in table '%s' on AirTable", conf.TableNames.Items)
		panic(errStr)
	}
	return items
}
