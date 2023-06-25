package entity

type Assert struct {
	ID           string
	Name         string
	MarketVolume int
}

func NewAssert(id string, name string, marketVolume int) *Assert {
	return &Assert{
		ID:           id,
		Name:         name,
		MarketVolume: marketVolume,
	}
}
