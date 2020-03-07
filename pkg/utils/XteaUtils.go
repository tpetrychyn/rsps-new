package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/xtea"
	"io/ioutil"
)

// making this because I can't decide if this should be global or injected somehow
var GlobalXteaDefs = make(map[uint16][]int32)

type XteaDefs = map[uint16][]int32

func XteaDecrypt(cipher *xtea.Cipher, src []byte) []byte {
	// pad to an even block size of 8
	if len(src) % xtea.BlockSize != 0 {
		pad := xtea.BlockSize - len(src) % xtea.BlockSize
		src = append(src, make([]byte, pad)...)
	}

	// allocate an empty slice with a capacity of src length
	res := make([]byte, 0, len(src))

	// iterate 8 bytes at a time decrypting them - one block at a time
	for i:=0;i<len(src)+xtea.BlockSize;i+=xtea.BlockSize {
		piece := src[i:i+xtea.BlockSize]
		dec := make([]byte, len(piece))
		cipher.Decrypt(dec, piece)
		res = append(res, dec...)
	}

	return res
}

func LoadXteas() (XteaDefs, error) {
	var xteaDefs = make(map[uint16][]int32)
	file, err := ioutil.ReadFile("xteas.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open xteas.json %w", err)
	}

	type XteaDef struct {
		Region uint16
		Keys []int32
	}

	var xteas []*XteaDef
	err = json.Unmarshal(file, &xteas)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xteas.json %w", err)
	}

	for _, v := range xteas {
		xteaDefs[v.Region] = v.Keys
	}

	return xteaDefs, nil
}