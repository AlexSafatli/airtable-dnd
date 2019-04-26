package main

import (
	"./store"
	"fmt"
	"github.com/evalphobia/go-config-loader"
)

const (
	confType = "toml"
	basePath = "config"
)

func main() {
	conf := config.NewConfig()
	if err := conf.LoadConfigs(basePath, confType); err != nil {
		panic(err)
	}
	conn, err := store.OpenConnection(conf.ValueString("database.api_key"), conf.ValueString("campaign.base_id"))
	if err != nil {
		panic(err)
	}
	characters := store.GetCharacters(conf.ValueString("characters.table_name"), conn)
	for i := range characters {
		fmt.Println(characters[i])
	}
}
