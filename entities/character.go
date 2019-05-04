package entities

type Character struct {
	Name            string
	Race            string
	Level           uint // == CR for Monsters
	Class           string
	HP              uint
	AC              uint
	Initiative      int
	PassivePer      uint `json:"Passive Per"`
	DexST           int  `json:"Dex ST"`
	Size            CharacterSize
	Type            string
	Speed           uint8
	Player          string
	Affiliated      bool
	Gender          string // {Male, Female}
	GenerationDelta int    `json:"Generation Delta"`
	YOBDelta        int    `json:"Year of Birth Delta"`
	YODDelta        int    `json:"Year of Death Delta"`
	MotherID        uint   `json:"Mother ID"`
	FatherID        uint   `json:"Father ID"`
	LocationID      uint   `json:"Location ID"`
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
