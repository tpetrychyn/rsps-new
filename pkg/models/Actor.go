package models

type Actor struct {
	Id            int
	Position      *Position
	NearbyPlayers []*Actor
	NearbyNpcs    []*Actor
	OutgoingQueue chan interface{}
}
