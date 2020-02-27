package main

import (
	"log"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet/outgoing"
	"time"
)

func main() {
	log.Printf("hello")

	p := game.NewPlayer()

	_ = outgoing.NewPlayerUpdatePacket(p.Actor)

	r := models.NewRegion()

	go func() {
		for {
			log.Printf("known %+v", p.Actor)
			popped := <- p.Actor.OutgoingQueue
			if m, ok := popped.(*models.Message); ok {
				log.Printf("popped message: %s", m.Body)
			}
		}
	}()

	for {
		p.AppendOutgoing()
		<- time.After(1 * time.Second)
		r.AddPlayer(p.Actor)
		<- time.After(1 * time.Second)
	}
}
