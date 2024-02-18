package config

import "github.com/evalphobia/go-config-loader"

const (
	confType        = "toml"
	basePath        = "config"
	apiKey          = "database.api_key"
	campaignBase    = "campaign.players_base_id"
	dmBase          = "campaign.dm_base_id"
	charactersTable = "characters.table_name"
	encountersTable = "encounters.table_name"
	itemsTable      = "items.table_name"
	npcsTable       = "npcs.table_name"
)

type Values struct {
	ApiKey       string
	CampaignBase string
	TableNames   struct {
		Characters string
		Encounters string
		Items      string
	}
}

func LoadConfigs() Values {
	var conf *config.Config
	conf = config.NewConfig()
	if err := conf.LoadConfigs(basePath, confType); err != nil {
		panic(err)
	}
	return Values{
		ApiKey:       conf.ValueString(apiKey),
		CampaignBase: conf.ValueString(campaignBase),
		TableNames: struct {
			Characters string
			Encounters string
			Items      string
		}{
			Characters: conf.ValueString(charactersTable),
			Encounters: conf.ValueString(encountersTable),
			Items:      conf.ValueString(itemsTable),
		},
	}
}
