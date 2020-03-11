package utils

import "osrs-cache-parser/pkg/cachestore"

var store = cachestore.NewStore()
func GetCache() *cachestore.Store {
	if store == nil {
		store = cachestore.NewStore()
	}
	return store
}
