package cli

import (
	"fmt"
	"github.com/AlexSafatli/airtable-dnd/config"
	"github.com/AlexSafatli/airtable-dnd/fetools"
	"github.com/AlexSafatli/airtable-dnd/lib"
	"github.com/fabioberger/airtable-go"
	"strings"
)

const coinWeight = 1 / 50.0

func manageItems(jsons []string, conf config.Values, conn *airtable.Client) {
	// Get party items
	var items = lib.GetAirtableItemsWithIDs(conf, conn)

	// TODO: call 5etools API microservice to get this data; for now, read JSON
	var itemMap = fetools.Get5etoolsItemMap(jsons)

	for i := range items {
		var name = strings.ToLower(items[i].Fields.Name)
		if (name == "party gold" || name == "gold") && items[i].Fields.Weight != coinWeight {
			// Set baseline for party gold and do not go further
			var updateGoldRecord = map[string]interface{}{
				"Appr. value": 1,          // a gold coin is worth 1 gp
				"Weight":      coinWeight, // a standard coin weights 1/50 lb
			}
			if err := lib.UpdateItemByID(items[i].AirtableID, updateGoldRecord, conf.TableNames.Items, conn); err != nil {
				panic(err)
			}
			fmt.Printf("Set weight baseline for party gold (name '%s') as weight %.2f lb.\n", items[i].Fields.Name, coinWeight)
			continue
		}
		var updateRecord = make(map[string]interface{}, 2)
		if v, ok := itemMap[name]; ok {
			updateRecord["Appr. value"] = v.Value / 10 // value is in SP
			updateRecord["Weight"] = v.Weight
			if err := lib.UpdateItemByID(items[i].AirtableID, updateRecord, conf.TableNames.Items, conn); err != nil {
				panic(err)
			}
			fmt.Printf("Updated item with name '%s' with value %.2f gp and weight %.2f lb.\n", v.Name, v.Value/10, v.Weight)
		}
	}
}
