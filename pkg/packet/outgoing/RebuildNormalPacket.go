package outgoing

import (
	"bytes"
	"encoding/binary"
	"log"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

const ChunkSize = 8
const ChunksPerRegion = 13

const MaxViewport = ChunkSize * ChunksPerRegion

type RebuildNormalPacket struct {
	Position *models.Position
}

func (r *RebuildNormalPacket) Build() []byte {
	rebuildNormalBuffer := utils.NewStream()

	rebuildNormalBuffer.WriteWordLEA(uint(r.Position.Z))
	rebuildNormalBuffer.WriteWordA(uint(r.Position.X))

	b := r.GetRegionXteas()
	rebuildNormalBuffer.Write(b)

	by := rebuildNormalBuffer.Flush()
	out := utils.NewStream()
	out.WriteByte(0)
	out.WriteWord(uint(len(by)))
	out.Write(by)

	return out.Flush()
}

func (r *RebuildNormalPacket) GetRegionXteas() []byte {
	lx := (r.Position.X - (MaxViewport >> 4)) >> 3
	rx := (r.Position.X + (MaxViewport >> 4)) >> 3
	lz := (r.Position.Z - (MaxViewport >> 4)) >> 3
	rz := (r.Position.Z + (MaxViewport >> 4)) >> 3

	buf := bytes.NewBuffer(make([]byte, 0, 2 + 4*10))

	forceSend := false
	if (r.Position.X / 8 == 48 || r.Position.X / 8 == 49) && r.Position.Z / 8 == 48 {
		forceSend = true
	}

	if r.Position.X / 8 == 48 && r.Position.Z / 8 == 48 {
		forceSend = true
	}

	// 386, 437
	log.Printf("posx %d posz %d", r.Position.X, r.Position.Z)

	count := 0
	buf.Write([]byte{0,0}) // make space for size short
	for x:=lx;x<=rx;x++ {
		for z:=lz;z<=rz;z++ {
			if !forceSend || z != 49 && z != 149 && z != 147 && x != 50 && (x != 49 || z != 47) {
				region := z + (x << 8)
				keys := utils.GlobalXteaDefs[region]
				binary.Write(buf, binary.BigEndian, keys)
				count++
			}
		}
	}

	by := buf.Bytes()
	by[0] = byte(count << 8)
	by[1] = byte(count & 0xFF)

	return by
}
