package counter_shard

import (
	"errors"
	"hash/fnv"
	"sync"
)

type CounterShard struct {
	Shards []ShardMap
	// Поля для хранилища (например, Redis, in-memory map, DB)
}

type ShardMap struct {
	mapping map[string]int
	rwMut   *sync.RWMutex
}

func New() *CounterShard {
	shards := make([]ShardMap, 10)
	for i := range shards {
		shards[i].rwMut = &sync.RWMutex{}
		shards[i].mapping = make(map[string]int)
	}

	return &CounterShard{
		shards,
	}
}

func (c *CounterShard) Increment(pageID string) error {
	if pageID == "" {
		return errors.New("pageID cannot be empty")
	}
	shardInex := hash(pageID, 10)
	shardMap := c.Shards[shardInex]

	shardMap.rwMut.Lock()
	shardMap.mapping[pageID]++
	shardMap.rwMut.Unlock()

	return nil
}

func (c *CounterShard) GetCount(pageID string) (int, error) {
	if pageID == "" {
		return 0, errors.New("pageID cannot be empty")
	}
	shardInex := hash(pageID, 10)
	shardMap := c.Shards[shardInex]
	shardMap.rwMut.RLock()
	defer shardMap.rwMut.RUnlock()
	count, exists := shardMap.mapping[pageID]
	if !exists {
		return 0, errors.New("pageID not found")
	}
	return count, nil
}

// хэширует и возвращает остаток
func hash(pageID string, numShards int) int {
	h := fnv.New32a()
	h.Write([]byte(pageID))
	hashValue := h.Sum32()
	return int(hashValue) % numShards
}
