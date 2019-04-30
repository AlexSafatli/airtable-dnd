package entities

type Character struct {
	Name       string
	Race       string
	Level      uint // == CR for Monsters
	Class      string
	HP         uint
	AC         uint
	Initiative int
	PassivePer uint `json:"Passive Per"`
	DexST      int  `json:"Dex ST"`
	Size       CharacterSize
	Type       string
	Speed      uint8
	Player     string
	Affiliated bool
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
