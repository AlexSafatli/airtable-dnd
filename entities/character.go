package entities

type Character struct {
	Name       string
	Race       string
	Level      uint
	Class      string
	HP         uint
	AC         uint
	Initiative int
	PassivePer uint `json:"Passive Per"`
	DexST      int  `json:"Dex ST"`
	Player     string
}
