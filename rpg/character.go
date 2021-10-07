package rpg

type Character struct {
	Name            string `json:",omitempty"`
	Race            string `json:",omitempty"`
	Level           uint   `json:",omitempty"`
	CR              uint   `json:",omitempty"`
	Class           string `json:",omitempty"`
	HP              uint
	AC              uint
	Initiative      int
	PassivePer      uint `json:"Passive Per"`
	DexST           int  `json:"Dex ST"`
	Size            CharacterSize
	Type            string `json:",omitempty"`
	Speed           uint8
	Player          string `json:",omitempty"`
	Affiliated      bool
	Gender          string `json:",omitempty"` // {Male, Female}
	GenerationDelta int    `json:"Generation Delta,omitempty"`
	YOBDelta        int    `json:"Year of Birth Delta,omitempty"`
	YODDelta        int    `json:"Year of Death Delta,omitempty"`
	MotherID        uint   `json:"Mother ID,omitempty"`
	FatherID        uint   `json:"Father ID,omitempty"`
	LocationID      uint   `json:"Location ID,omitempty"`
}

type CharacterSize uint8

const (
	SizeTiny       CharacterSize = 0
	SizeSmall      CharacterSize = 1
	SizeMedium     CharacterSize = 2
	SizeLarge      CharacterSize = 3
	SizeHuge       CharacterSize = 4
	SizeGargantuan CharacterSize = 5
)
