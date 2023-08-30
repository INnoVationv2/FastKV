package lru

import (
	"fmt"
	"math/rand"
	"testing"
)

func GetHead(l *LRU) *Node {
	return l.head.next
}

func GetTail(l *LRU) *Node {
	return l.head.prev
}

// 检测是否能持续淘汰，保证LRU不超出最大容量
func TestLRU_Add(t *testing.T) {
	sz := 100
	lru := newLRU(sz)
	strSet := GetRandomString(10000)
	for i := 1; i <= len(strSet); i++ {
		pair := strSet[i-1]
		lru.Add(pair.key, pair.value)

		lruLen := lru.length()
		if i < sz {
			assert(lru.length() == i, fmt.Sprintf("%d != %d", lruLen, i))
		} else {
			assert(lruLen == sz, fmt.Sprintf("%d != %d", lruLen, sz))
		}
	}
}

// 检测随机访问元素，元素能否移动到队头
func TestLRU_Get(t *testing.T) {
	sz := 1000
	lru := newLRU(sz)
	strSet := GetRandomString(1000)
	for i := 0; i < len(strSet); i++ {
		pair := strSet[i]
		lru.Add(pair.key, pair.value)
	}

	for i := 0; i < 1000; i++ {
		val := rand.Intn(1000)
		pair := strSet[val]
		lru.Get(pair.key)
		assert(GetHead(lru).value == pair.value, "Front Val Not Equal to"+pair.value)
	}
}

// 检测淘汰的值是否正确, 每次淘汰，末尾的值会改变
func TestLRU_Del(t *testing.T) {
	sz := 10
	lru := newLRU(sz)
	strSet := GetRandomString(1000)
	for i := 0; i < len(strSet); i++ {
		pair := strSet[i]
		lru.Add(pair.key, pair.value)
		if i < sz {
			continue
		}
		assert(GetTail(lru).value == strSet[i-9].value, "Tail Node not except value.")
	}
}

func assert(condition bool, errorMsg string) {
	if !condition {
		panic(errorMsg)
	}
}
