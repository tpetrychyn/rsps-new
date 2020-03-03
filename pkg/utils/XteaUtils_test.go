package utils

import (
	"golang.org/x/crypto/xtea"
	"log"
	"testing"
)

func Test_Xtea(t *testing.T) {
	xteaKeys := make([]int32, 4)
	xteaKeys[0] = 1157681262
	xteaKeys[1] = -2056262709
	xteaKeys[2] = 1733026640
	xteaKeys[3] = 1210958896

	xteaKey := make([]byte, 16)
	for i := 0; i < len(xteaKeys); i++ {
		j := i << 2
		xteaKey[j] = byte(xteaKeys[i]>>24)
		xteaKey[j+1] = byte(xteaKeys[i]>>16)
		xteaKey[j+2] = byte(xteaKeys[i]>>8)
		xteaKey[j+3] = byte(xteaKeys[i])
	}

	log.Printf("xteaKey %v", xteaKey)

	xteaCipher, err := xtea.NewCipher(xteaKey)
	if err != nil {
		panic(err)
	}

	buf := []byte{218, 53, 225, 2, 140, 220, 42, 126, 241, 10, 150, 251, 104, 22, 167, 61, 191, 192, 236, 26, 22, 47, 168, 131, 33, 141, 0, 134, 201, 130, 126, 105, 90, 13, 202, 184, 213, 201, 81, 37, 51, 179, 57, 62, 0, 147, 211, 47, 72, 29, 112, 22, 228, 82, 99, 249, 132, 249, 144, 75, 194, 58, 98, 213, 50, 41, 74, 175, 96, 109, 60, 68, 49, 217, 226, 187, 251, 65, 63, 61, 124, 46, 116, 231, 110, 98, 79, 226, 212, 8, 8, 173, 158, 175, 65, 225, 83, 101, 18, 121, 81, 198, 170, 86, 191, 93, 11, 192, 58, 165, 159, 110, 15, 253, 69, 201, 110, 72, 246, 20, 177, 212, 233, 167, 52, 181, 161, 207, 171, 114, 184, 79, 174, 90, 145, 46, 171, 114, 184, 79, 174, 90, 145, 46, 171, 114, 184, 79, 174, 90, 145, 46, 171, 114, 184, 79, 174, 90, 145, 46, 72, 171, 165, 9, 67, 63, 102, 128, 141, 37, 177, 15, 247, 80, 209, 114, 242, 129, 22, 249, 89, 167, 170, 201, 236, 14, 182, 129, 241, 103, 36, 144, 75, 159, 111, 39, 158, 134, 51, 210, 246, 95, 10, 138, 33, 10, 99, 83, 47, 153, 204, 238, 28, 22, 208, 173, 241, 119, 109, 176, 220, 147, 176, 14, 62, 73, 172, 45, 120, 131, 78, 120, 195, 199, 9, 0, 148, 24, 238, 12, 167, 100, 219, 95, 130, 182, 152, 207, 107, 98, 61}

	if len(buf) % 8 != 0 {
		log.Printf("len buf %d", len(buf))
		pad := 8 - len(buf) % 8
		buf = append(buf, make([]byte, pad)...)
	}

	res := make([]byte, 0, len(buf))

	for i:=0;i<len(buf)+8;i+=8 {
		piece := buf[i:i+8]
		dec := make([]byte, len(piece))
		xteaCipher.Decrypt(dec, piece)
		res = append(res, dec...)
	}

	if string(res[:20]) != "tpetrychyn@gmail.com" {
		t.Fatalf("expected tpetrychyn@gmail.com got %s", string(res[:20]))
	}
}
