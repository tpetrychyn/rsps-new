package main

import (
	"bufio"
	"encoding/binary"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/archives"
	"github.com/tpetrychyn/rsps-comm-test/internal/game"
	"github.com/tpetrychyn/rsps-comm-test/internal/game/entities"
	rsNet "github.com/tpetrychyn/rsps-comm-test/internal/net"
	"github.com/tpetrychyn/rsps-comm-test/pkg/models"
	"github.com/tpetrychyn/rsps-comm-test/pkg/packet/outgoing"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
	"log"
	"net"
	"time"
)

const revision = 181
const port = "43594"

func main() {
	xteaDefs, err := utils.LoadXteas()
	if err != nil {
		panic(err)
	}
	utils.GlobalXteaDefs = xteaDefs

	cacheStore := utils.GetCache()
	js5Handler := rsNet.JS5Handler{CacheStore: cacheStore}
	loginHandler := rsNet.LoginHandler{CacheStore: cacheStore}

	objArchive := archives.NewObjectArchive(cacheStore)
	utils.GetDefinitions().SetObjects(objArchive.LoadObjectDefs())

	world := &game.World{Chunks: make(map[string]*models.Chunk)}
	players := make(map[int]*rsNet.Client)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

	log.Printf("Listening on port %s", port)

	go func() {
		for {
			for _, client := range players {
				if client == nil {
					continue
				}
				client.Player.Pretick()
				client.Player.Tick()
			}

			for _, client := range players {
				if client == nil {
					continue
				}
				client.Player.PostTick()
			}

			<-time.After(600 * time.Millisecond)
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		reader := bufio.NewReader(connection)

		var requestType byte
		binary.Read(reader, binary.BigEndian, &requestType)

		if requestType == 15 {
			var gameVersion int32
			binary.Read(reader, binary.BigEndian, &gameVersion)
			if gameVersion != revision {
				connection.Write([]byte{6}) // out of date
				continue
			}

			connection.Write([]byte{0})
			go func() {
				// continuously loop through reading js5 request bytes for this socket until error
				for js5Handler.HandleRequest(connection, reader) {
				}
				connection.Close()
			}()
		}

		if requestType == 14 {
			encryptor, decryptor := loginHandler.HandleRequest(connection, reader)
			if encryptor == nil || decryptor == nil {
				connection.Close()
				continue
			}

			player := entities.NewPlayer(world)
			client := rsNet.NewClient(connection, encryptor, decryptor, player)
			player.OutgoingQueue = client.DownstreamQueue

			player.Actor.Teleport(&models.Tile{X: 3221, Y: 3222})

			client.DownstreamQueue <- &outgoing.RebuildLoginPacket{Position: player.Actor.Movement.Position}

			client.DownstreamQueue <- &outgoing.IfOpenTopPacket{} // main screen interface?

			client.Player.Movement.IsRunning = true

			players[player.Id] = client
		}
	}
}
