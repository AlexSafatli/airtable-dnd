package fetools

type Spell struct {
	Name   string `json:"name"`
	Level  int    `json:"level"`
	School string `json:"school"`
	//Time SpellTime `json:"time"`
	Range       SpellRange    `json:"range"`
	Classes     SpellClasses  `json:"classes"`
	Source      string        `json:"source"`
	Entries     []interface{} `json:"entries"`
	Page        int           `json:"page"`
	DamageTypes []string      `json:"damageInflict"`
}

type SpellTime struct {
	Number int    `json:"number"`
	Unit   string `json:"unit"`
}

type SpellRange struct {
	Type     string        `json:"type"`
	Distance SpellDistance `json:"disance"`
}

type SpellDistance struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type SpellClasses struct {
	ClassList []Class `json:"fromClassList"`
}

type Class struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

type Monster struct {
	Name   string    `json:"name"`
	Size   string    `json:"size"`
	Source string    `json:"source"`
	HP     MonsterHP `json:"hp"`
	Dex    int       `json:"dex"`
}

type MonsterType struct {
	Type string   `json:"type"`
	Tags []string `json:"tags"`
}

type MonsterHP struct {
	Average int    `json:"average"`
	Formula string `json:"formula"`
}

type Item struct {
	Name       string        `json:"name"`
	Source     string        `json:"source"`
	Page       int           `json:"page"`
	Rarity     string        `json:"rarity"`
	Wondrous   bool          `json:"wondrous"`
	Entries    []interface{} `json:"entries"`
	Weight     float32       `json:"weight"`
	BaseItem   string        `json:"baseItem"`
	Type       string        `json:"type"`
	Property   []interface{} `json:"property"`
	DamageType string        `json:"dmgType"`
	Tier       string        `json:"tier"`
	LootTables []interface{} `json:"lootTables"`
	Value      float32       `json:"value"`
}
