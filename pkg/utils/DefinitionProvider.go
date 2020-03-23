package utils

import "github.com/tpetrychyn/osrs-cache-parser/pkg/models"

type DefinitionContainer struct {
	Objects []*models.ObjectDef
}

var definitions = &DefinitionContainer{}

func GetDefinitions() *DefinitionContainer {
	return definitions
}

func (d *DefinitionContainer) SetObjects(objs []*models.ObjectDef) {
	definitions.Objects = objs
}
