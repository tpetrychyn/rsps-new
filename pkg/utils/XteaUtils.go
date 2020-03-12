package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// making this because I can't decide if this should be global or injected somehow
var GlobalXteaDefs = make(map[uint16][]int32)

type XteaDefs = map[uint16][]int32

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