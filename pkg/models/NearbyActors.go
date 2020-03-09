package models

import "sync"

type NearbyActors struct {
	mut *sync.Mutex
	array []*Actor
}

func NewNearbyActors() *NearbyActors {
	return &NearbyActors{
		mut:   new(sync.Mutex),
		array: make([]*Actor, 255),
	}
}

func (n *NearbyActors) Get() []*Actor {
	n.mut.Lock()
	defer n.mut.Unlock()
	return n.array
}

func (n *NearbyActors) Set(players []*Actor) {
	n.mut.Lock()
	n.array = players
	n.mut.Unlock()
}