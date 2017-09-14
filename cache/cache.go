package cache

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/rikardNL/tegola"
)

// Cache is a contract for tile cache to be implemented by specializations of caches
type Cache interface {
	Set(tile tegola.Tile, layerName string, mapName string, value []byte) error
	Get(tile tegola.Tile, layerName string, mapName string) ([]byte, error)
}

// RedisCache is an implementation of Cache using Redis as storage backend
type RedisCache struct {
	RedisClient redis.Client
	TTL         time.Duration
}

// keyFor consistently converts a tile to the key it will be associated with in the redis cache
func (r RedisCache) keyFor(tile tegola.Tile, layerName string, mapName string) string {
	return fmt.Sprintf("%s/%s/%d/%d/%d", mapName, layerName, tile.Z, tile.X, tile.Y)
}

// Set the tile to the supplied byte array
func (r RedisCache) Set(tile tegola.Tile, layerName string, mapName string, value []byte) error {
	key := r.keyFor(tile, layerName, mapName)
	return r.RedisClient.Set(key, value, r.TTL).Err()
}

// Get the byte array stored for the tile
func (r RedisCache) Get(tile tegola.Tile, layerName string, mapName string) ([]byte, error) {
	key := r.keyFor(tile, layerName, mapName)
	return r.RedisClient.Get(key).Bytes()
}
