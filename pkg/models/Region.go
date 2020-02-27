package models

import (
	"log"
	"sync"
)

type Region struct {
	Players *sync.Map
	Npcs    *sync.Map
}

func NewRegion() *Region {
	return &Region{
		Players: &sync.Map{},
		Npcs:    &sync.Map{},
	}
}

func (r *Region) AddPlayer(player *Actor) {
	r.Players.Store(player.Id, player)
	r.Players.Range(func(key, value interface{}) bool {
		if p, ok := value.(*Actor); ok {
			log.Printf("found %+v", p)
			p.OutgoingQueue <- &Message{Body:"region enter"}
		}
		return true
	})
}
