package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"osrs-cache-parser/pkg/cachestore"
)

type JS5Handler struct {
	CacheStore *cachestore.Store
}

func (j *JS5Handler) HandleRequest(connection net.Conn, reader *bufio.Reader) bool {
	var opcode byte
	err := binary.Read(reader, binary.BigEndian, &opcode)
	if err != nil {
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
				for _, v := range j.CacheStore.Indexes {
					binary.Write(dataBuffer, binary.BigEndian, v.Crc)
					binary.Write(dataBuffer, binary.BigEndian, v.Revision)
				}
				data.Write([]byte{0}) // no compression
				binary.Write(data, binary.BigEndian, int32(dataBuffer.Len()))
				binary.Write(data, binary.BigEndian, dataBuffer.Bytes())
			} else {
				r := j.CacheStore.ReadIndex(int(archive))
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
			i := j.CacheStore.FindIndex(int(index))
			a := i.Archives[archive]
			data := j.CacheStore.LoadArchive(a)

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