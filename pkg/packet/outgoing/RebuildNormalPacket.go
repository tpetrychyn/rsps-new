package outgoing

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"log"
)

const ChunkSize = 8
const ChunksPerRegion = 13

const MaxViewport = ChunkSize * ChunksPerRegion

type RebuildNormalPacket struct {
	X uint32
	Z uint32
}

func (r *RebuildNormalPacket) Write(writer *bufio.Writer) []byte {
	lx := (r.X - (MaxViewport >> 4)) >> 3
	rx := (r.X + (MaxViewport >> 4)) >> 3
	lz := (r.Z - (MaxViewport >> 4)) >> 3
	rz := (r.Z + (MaxViewport >> 4)) >> 3

	//lx = 47
	//rx = 49
	//lz = 53
	//rz = 55

	buf := bytes.NewBuffer(make([]byte, 0, 2 + 4*10))

	forceSend := false
	if (r.X / 8 == 48 || r.X / 8 == 49) && r.Z / 8 == 48 {
		forceSend = true
	}

	if r.X / 8 == 48 && r.Z / 8 == 48 {
		forceSend = true
	}

	count := 0
	buf.Write([]byte{0,0}) // make space for size short
	for x:=lx;x<=rx;x++ {
		for z:=lz;z<=rz;z++ {
			if !forceSend || z != 49 && z != 149 && z != 147 && x != 50 && (x != 49 || z != 47) {
				//region := z + (x << 8)
				// TODO: load xtea.json, get xtea[region]
				keys := []int32{-1621279394, 1091446180, -341021727, -599622964}
				binary.Write(buf, binary.BigEndian, keys)
				count++
			}
		}
	}

	by := buf.Bytes()
	by[0] = byte(count << 8) // should be 9
	by[1] = byte(count & 0xFF)

	log.Printf("length %d buf bytes %+v", len(buf.Bytes()), buf.Bytes()) // should be 146 length

	//
	return []byte{0, 9, 159, 93, 61, 94, 65, 14, 37, 164, 235, 172, 107, 225, 220, 66, 122, 204, 181, 20, 30, 174, 148, 75, 142, 166, 201, 249, 182, 89, 235, 183, 9, 102, 59, 93, 251, 245, 30, 206, 103, 243, 138, 169, 38, 254, 175, 55, 60, 156, 158, 62, 97, 0, 47, 8, 111, 152, 44, 130, 101, 124, 128, 176, 22, 157, 181, 248, 165, 66, 43, 50, 93, 159, 247, 63, 20, 161, 68, 132, 9, 124, 179, 58, 54, 81, 32, 72, 250, 232, 102, 193, 182, 39, 198, 236, 166, 38, 188, 198, 159, 213, 194, 63, 119, 240, 228, 236, 152, 30, 9, 37, 248, 182, 14, 67, 173, 76, 106, 98, 130, 161, 42, 31, 130, 175, 249, 169, 147, 205, 223, 50, 137, 71, 39, 75, 149, 77, 24, 154, 125, 28, 109, 87, 208, 17}
	//writer.Write(by)
}
