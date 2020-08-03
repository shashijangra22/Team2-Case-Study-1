package Models

type Item struct {
	Name     string
	Price    float32
}


type Rest struct {
	ID    int64
	Name string
	Availability bool
	Items   []Item
	Rating   float32
	Category  string
}
