package counter_shard

import (
	"errors"
	"hash/fnv"
	"sync"
)

type CounterShard struct {
	Shards        []ShardMap
	incrementChan chan string
	// Поля для хранилища (например, Redis, in-memory map, DB)
}

type ShardMap struct {
	mapping map[string]int
	rwMut   *sync.RWMutex
}

func New() *CounterShard {
	shardChan := make(chan string, 100_000)
	shards := make([]ShardMap, 10)
	for i := range shards {
		shards[i].rwMut = &sync.RWMutex{}
		shards[i].mapping = make(map[string]int)
	}

	res := &CounterShard{
		shards,
		shardChan,
	}
	go res.incrementShard(shardChan)
	return res
}

func (c *CounterShard) Increment(pageID string) error {
	if pageID == "" {
		return errors.New("pageID cannot be empty")
	}

	c.incrementChan <- pageID

	return nil
}

// mini прога для записи в мапу, в которую стекаются все сообщения из обработчиков
func (c *CounterShard) incrementShard(ch chan string) {
	for pageID := range ch {
		shardInex := hash(pageID, 10)
		shardMap := c.Shards[shardInex]
		shardMap.rwMut.Lock()
		shardMap.mapping[pageID]++
		shardMap.rwMut.Unlock()
	}
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
