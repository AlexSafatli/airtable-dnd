package main

import "github.com/evalphobia/go-config-loader"

const (
	confType        = "toml"
	basePath        = "config"
	apiKey          = "database.api_key"
	campaignBase    = "campaign.players_base_id"
	dmBase          = "campaign.dm_base_id"
	charactersTable = "characters.table_name"
	encountersTable = "encounters.table_name"
	npcsTable       = "npcs.table_name"
)

type configValues struct {
	ApiKey       string
	CampaignBase string
	TableNames   struct {
		Characters string
		Encounters string
	}
}

func loadConfigs() configValues {
	var conf *config.Config
	conf = config.NewConfig()
	if err := conf.LoadConfigs(basePath, confType); err != nil {
		panic(err)
	}
	return configValues{
		ApiKey:       conf.ValueString(apiKey),
		CampaignBase: conf.ValueString(campaignBase),
		TableNames: struct {
			Characters string
			Encounters string
		}{
			Characters: conf.ValueString(charactersTable),
			Encounters: conf.ValueString(encountersTable),
		},
	}
}
