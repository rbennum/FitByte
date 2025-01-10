package cache

import (
	"encoding/json"
	"github.com/dgraph-io/ristretto/v2"
	"log"
	"time"
)

const (
	MaxCacheSize = 256 << 20 //  256 MB
	DefaultTtl   = 30 * time.Minute

	CacheAuthEmailToToken      = "auth:%s"
	CacheUserIdToProfile       = "user:%s"
	CacheEmployeesWithParams   = "employees:%s"
	CacheDepartmentsWithParams = "departments:%s"
)

var Cache *ristretto.Cache[string, string]

func initialize() {
	var err error

	Cache, err = ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e6, // 1 million counters for frequency tracking
		MaxCost:     MaxCacheSize,
		BufferItems: 64, // 64 keys per Get buffer
	})
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}
}

func Set(key string, value string) {
	cost := int64(len(key) + len(value))
	SetWithCost(key, value, cost)
}

func SetAsMap(key string, value map[string]string) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	cost := int64(len(key) + len(data))

	SetWithCost(key, string(data), cost)
}

func SetAsMapArrayWithTtlAndCostMultiplier(
	key string,
	value []map[string]string,
	costMultiplier int,
	ttl time.Duration,
) {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	totalCost := int64(len(key))
	for _, v := range value {
		for k, str := range v {
			stringCost := int64(len(k)) + int64(len(str))
			totalCost += stringCost
		}
	}

	Cache.SetWithTTL(key, string(data), totalCost*int64(costMultiplier), ttl)
}

func SetWithCost(key string, value string, cost int64) {
	Cache.SetWithTTL(key, value, cost, DefaultTtl)
}

func Get(key string) (string, bool) {
	return Cache.Get(key)
}

func GetAsMap(key string) (map[string]string, bool) {
	val, found := Cache.Get(key)
	if !found {
		return nil, false
	}

	var result map[string]string
	err := json.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}

	return result, true
}

func GetAsMapArray(key string) ([]map[string]string, bool) {
	val, found := Cache.Get(key)
	if !found {
		return nil, false
	}

	var result []map[string]string
	err := json.Unmarshal([]byte(val), &result)
	if err != nil {
		panic(err)
	}

	return result, true
}

func Delete(key string) {
	Cache.Del(key)
}
