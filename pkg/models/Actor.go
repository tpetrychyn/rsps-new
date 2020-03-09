package models

type Actor struct {
	Id            int
	Movement      *Movement
	NearbyPlayers *NearbyActors
	NearbyNpcs    *NearbyActors
}

func NewActor() *Actor {
	return &Actor {
		Id:            0,
		Movement:      NewMovement(),
		NearbyPlayers: NewNearbyActors(),
		NearbyNpcs:    NewNearbyActors(),
	}
}
