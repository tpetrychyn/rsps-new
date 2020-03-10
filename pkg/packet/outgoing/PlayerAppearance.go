package outgoing

import (
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/utils"
)

type PlayerAppearance struct {
	Target *models.Actor
}

func (p *PlayerAppearance) Build() []byte {
	stream := utils.NewStream()

	stream.WriteByte(0)   // gender
	stream.WriteByte(255) // skull icon
	stream.WriteByte(255) // prayer icon

	arms, hair, _ := 6, 8, 11
	for i:=0;i<12;i++ {
		if i == arms {
			item := p.Target.Equipment.Items[4]
			if item != nil {
				stream.WriteByte(0)
				continue
			}
		}
		if i == hair {
			item := p.Target.Equipment.Items[0]
			if item != nil {
				stream.WriteByte(0)
				continue
			}
		}

		item := p.Target.Equipment.Items[i]
		if item != nil {
			stream.WriteWord(0x200 + uint(item.Id))
		} else {
			if translation[i] == -1 {
				stream.WriteByte(0)
			} else {
				stream.WriteWord(0x100 + defaultLooks[translation[i]])
			}
		}
	}

	for _, c := range defaultColors {
		stream.WriteByte(c)
	}

	weapon := p.Target.Equipment.Items[3] // weapon slot 3
	if weapon != nil {
		// TODO: Change the animations
	}

	for _, a := range animations {
		stream.WriteWord(a)
	}

	stream.WriteString("tpetrychyn@gmail.com")
	stream.WriteByte(3)
	stream.WriteWord(0)
	stream.WriteByte(0)

	by := stream.Flush()
	size := len(by)
	for k := range by {
		by[k] += 128
	}
	return append([]byte{byte(128-size)}, by...)
}

var translation = []int{-1, -1, -1, -1, 2, -1, 3, 5, 0, 4, 6, 1}

var animations = []uint{808, 823, 819, 820, 821, 822, 824}
var defaultLooks = []uint{9, 14, 109, 26, 33, 36, 42}
var defaultColors = []byte{0, 3, 2, 0, 0}

// [19 0 73 192 127 244 5 0 0 62 128 127 127 128 128 128 128 129 237 128 129 154 129 164 129 137 129 161 129 170 129 142 128 131 130 128 128 131 168 131 183 131 179 131 180 131 181 131 182 131 184 244 240 229 244 242 249 227 232 249 238 192 231 237 225 233 236 174 227 239 237 128 131 128 128 128 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]

// [79 0 73 240 0 0 127 244 1 62 128 127 127 128 128 128 128 129 237 128 129 154 129 164 129 137 129 161 129 170 129 142 128 131 130 128 128 131 168 131 183 131 179 131 180 131 181 131 182 131 184 244 240 229 244 242 249 227 232 249 238 192 231 237 225 233 236 174 227 239 237 128 131 128 128 128]
