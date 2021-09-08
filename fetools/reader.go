package fetools

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

type SpellList struct {
	Spells []Spell `json:"spell"`
}

type MonsterList struct {
	Monsters []Monster `json:"monster"`
}

type ItemList struct {
	Items     []Item `json:"item"`
	BaseItems []Item `json:"baseitem"`
}

func Get5etoolsSpells(path string) []Spell {
	var spellList SpellList
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &spellList)
	if err != nil {
		panic(err)
	}
	return spellList.Spells
}

func Get5etoolsMonsters(path string) []Monster {
	var monsterList MonsterList
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &monsterList)
	if err != nil {
		panic(err)
	}
	return monsterList.Monsters
}

func Get5etoolsItems(paths []string) []Item {
	var itemList ItemList
	for i := range paths {
		var items ItemList
		dat, err := ioutil.ReadFile(paths[i])
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(dat, &items)
		if err != nil {
			panic(err)
		}
		if len(items.Items) > 0 {
			itemList.Items = append(itemList.Items, items.Items...)
		}
		if len(items.BaseItems) > 0 {
			itemList.Items = append(itemList.Items, items.BaseItems...)
		}
	}
	return itemList.Items
}

func Get5etoolsItemMap(paths []string) map[string]Item {
	var itemMap map[string]Item
	var items []Item
	items = Get5etoolsItems(paths)
	itemMap = make(map[string]Item)
	for i := range items {
		itemMap[strings.ToLower(items[i].Name)] = items[i]
	}
	return itemMap
}
