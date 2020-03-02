package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"osrs-cache-parser/pkg/cachestore"
)

var store = cachestore.NewStore()

func main() {

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
			if gameVersion == 188 {
				connection.Write([]byte{0})
			}
			go func() {for handleJS5Requests(connection, reader) {}}()
		} else {
			log.Printf("requestType: %d", requestType)
		}
	}
}

func handleJS5Requests(connection net.Conn, reader *bufio.Reader) bool {
	var opcode byte
	err := binary.Read(reader, binary.BigEndian, &opcode)
	if err != nil {
		connection.Close()
		return false
	}

	if opcode == 3 || opcode == 2 || opcode == 6 { // skip loading screen
		reader.Discard(3)
		binary.Read(reader, binary.BigEndian, &opcode)
	}

	if opcode == 0 || opcode == 1 {
		var index uint8
		var archive uint16

		binary.Read(reader, binary.BigEndian, &index)
		binary.Read(reader, binary.BigEndian, &archive)

		if index == 255 {
			data := new(bytes.Buffer)
			if archive == 255 {
				dataBuffer := new(bytes.Buffer)
				for _, v := range store.Indexes {
					binary.Write(dataBuffer, binary.BigEndian, v.Crc)
					binary.Write(dataBuffer, binary.BigEndian, v.Revision)
				}
				data.Write([]byte{0}) // no compression
				binary.Write(data, binary.BigEndian, int32(dataBuffer.Len()))
				binary.Write(data, binary.BigEndian, dataBuffer.Bytes())
			} else {
				r := store.ReadIndex(int(archive))
				data.Write(r)
			}

			outBuffer := new(bytes.Buffer)
			binary.Write(outBuffer, binary.BigEndian, index)
			binary.Write(outBuffer, binary.BigEndian, archive)
			writeIndex := 3 // 1 byte for index, 2 bytes for archive so we've written 3 so far
			for _, v := range data.Bytes() {
				if writeIndex%512 == 0 {
					binary.Write(outBuffer, binary.BigEndian, byte(255))
					writeIndex++
				}
				binary.Write(outBuffer, binary.BigEndian, v)
				writeIndex++
			}

			connection.Write(outBuffer.Bytes())
		} else {
			i := store.FindIndex(int(index))
			a := i.Archives[archive]
			data := store.LoadArchive(a)

			if a == nil || data == nil {
				log.Printf("nil data index %d", index)
			}

			if data != nil {
				compression := data[0]
				length := int(data[1])<<24 | (int(data[2])&0xFF)<<16 | (int(data[3])&0xFF)<<8 | int(data[4])&0xFF
				expectedLength := length
				if compression == 0 {
					expectedLength += 5
				} else {
					expectedLength += 9
				}
				if len(data)-expectedLength == 2 {
					data = data[:len(data)-2]
				}

				outBuffer := new(bytes.Buffer)
				binary.Write(outBuffer, binary.BigEndian, index)
				binary.Write(outBuffer, binary.BigEndian, archive)
				writeIndex := 3 // 1 byte for index, 2 bytes for archive so we've written 3 so far
				for _, v := range data {
					if writeIndex%512 == 0 {
						binary.Write(outBuffer, binary.BigEndian, byte(255))
						writeIndex++
					}
					binary.Write(outBuffer, binary.BigEndian, v)
					writeIndex++
				}
				connection.Write(outBuffer.Bytes())
			}
		}
	}

	return true
}
