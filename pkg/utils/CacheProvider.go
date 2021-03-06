package utils

import "github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"

var cache = cachestore.NewStore("./cache")

func GetCache() *cachestore.Store {
	if cache == nil {
		cache = cachestore.NewStore("./cache")
	}

	return cache
}
