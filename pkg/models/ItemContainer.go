package models

type Item struct {
	Id int
	Amount int
}

type ItemContainer struct {
	Capacity uint
	Items []*Item
}

func NewItemContainer(capacity uint) *ItemContainer {
	return &ItemContainer{
		Capacity: capacity,
		Items:    make([]*Item, capacity),
	}
}
