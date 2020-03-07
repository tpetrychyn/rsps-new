package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"osrs-cache-parser/pkg/cachestore"
	rsNet "rsps-comm-test/internal/net"
	"rsps-comm-test/pkg/utils"
)

func main() {
	xteaDefs, err := utils.LoadXteas()
	if err != nil {
		panic(err)
	}
	utils.GlobalXteaDefs = xteaDefs

	cacheStore := cachestore.NewStore()

	js5Handler := rsNet.JS5Handler{CacheStore: cacheStore}
	loginHandler := rsNet.LoginHandler{CacheStore: cacheStore}

	listener, err := net.Listen("tcp", ":43594")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		reader := bufio.NewReader(connection)

		var requestType byte
		binary.Read(reader, binary.BigEndian, &requestType)

		log.Printf("request type %d", requestType)

		if requestType == 15 {
			var gameVersion int32
			binary.Read(reader, binary.BigEndian, &gameVersion)
			if gameVersion == 181 {
				connection.Write([]byte{0})
			}
			go func() {
				// continuously loop through reading js5 request bytes for this socket until error
				for js5Handler.HandleRequest(connection, reader) {}
				connection.Close()
			}()
		}

		if requestType == 14 {
			go func () {
				success := loginHandler.HandleRequest(connection, reader)
				if !success {
					connection.Close()
				}

				// otherwise we have a logged in player
			}()
		}
	}
}
