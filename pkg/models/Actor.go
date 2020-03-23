package models

type Actor struct {
	Id            int
	Movement      *Movement
	NearbyPlayers *NearbyActors
	NearbyNpcs    *NearbyActors
	UpdateMask    UpdateMask
	Equipment     *ItemContainer
}

func NewActor() *Actor {
	return &Actor{
		Id:            0,
		Movement:      NewMovement(),
		NearbyPlayers: NewNearbyActors(),
		NearbyNpcs:    NewNearbyActors(),
		Equipment:     NewItemContainer(14),
	}
}

func (a *Actor) Teleport(tile *Tile) {
	a.Movement.Position = tile
	a.UpdateMask.Movement = true
	a.Movement.Teleported = true
}