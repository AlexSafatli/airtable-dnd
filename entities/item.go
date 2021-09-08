package entities

type Item struct {
	Name        string
	Quantity    uint
	Weight      float32
	ApproxValue float32 `json:"Appr. value"`
	//Type        string
}
