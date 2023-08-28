package lru

import (
	"FastKV/cache"
	"math/rand"
	"strconv"
	"testing"
)

func TestLRU_Add(t *testing.T) {
	set := make(cache.Set)
	sz := 5
	loop := 10000
	lru := newLRU(sz)
	for i := 0; i < loop; i++ {
		val := strconv.Itoa(rand.Intn(100))
		K := "Key" + val
		V := "Val" + val
		lru.Add(K, V)
		set[K] = nil

		lruLen := int(lru.length())
		if len(set) >= sz {
			assert(lruLen == sz, "Not Equal", lruLen, sz)
		} else {
			assert(lruLen == len(set), "Not Equal", lruLen, sz)
		}
	}
}

func assert(condition bool, message string, x int, y int) {
	if !condition {
		print(x, " ", y)
		panic(message)
	}
}
