package rpg

import "sort"

type Encounter struct {
	Name       string
	Session    uint
	Location   string
	Level      uint
	Room       uint
	NumPlayers uint `json:"Number of Players"`
	XP         uint
}

type EncounterState struct {
	Round           uint
	Turn            uint
	InitiativeOrder []*Character
}

type InitiativeValue struct {
	Character  *Character
	Initiative int
}

type InitiativeList []InitiativeValue

func (il InitiativeList) Len() int { return len(il) }
func (il InitiativeList) Less(i, j int) bool {
	return il[i].Initiative < il[j].Initiative || (il[i].Initiative == il[j].Initiative && il[i].Character.DexST < il[j].Character.DexST)
}
func (il InitiativeList) Swap(i, j int) { il[i], il[j] = il[j], il[i] }

func (il InitiativeList) Characters() []*Character {
	var characters []*Character
	characters = make([]*Character, len(il))
	for i, init := range il {
		characters[i] = init.Character
	}
	return characters
}

func NewEncounter(session uint, loc string, level, room, numPlayers uint) *Encounter {
	return &Encounter{Session: session, Location: loc, Level: level, Room: room, NumPlayers: numPlayers}
}

func RankInitiatives(initiatives map[*Character]int) InitiativeList {
	il := make(InitiativeList, len(initiatives))
	i := 0
	for k, v := range initiatives {
		il[i] = InitiativeValue{k, v}
		i++
	}
	sort.Sort(sort.Reverse(il))
	return il
}
