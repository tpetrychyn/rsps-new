package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"osrs-cache-parser/pkg/cachestore"
	rsNet "rsps-comm-test/internal/net"
	"rsps-comm-test/pkg/utils"
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

	cacheStore := cachestore.NewStore()

	js5Handler := rsNet.JS5Handler{CacheStore: cacheStore}
	loginHandler := rsNet.LoginHandler{CacheStore: cacheStore}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

	log.Printf("Listening on port %s", port)

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
				for js5Handler.HandleRequest(connection, reader) {}
				connection.Close()
			}()
		}

		if requestType == 14 {
			go func () {
				client := loginHandler.HandleRequest(connection, reader)
				if client == nil {
					connection.Close()
					return
				}

				for {
					client.Player.Tick()
					//client.EnqueueOutgoing(outgoing.NewPlayerUpdatePacket(client.Player.Actor))
					<- time.After(600 * time.Millisecond)
				}
			}()
		}
	}
}
