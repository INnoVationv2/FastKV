package O_LogN

import (
	"log"
	"math/rand"
	"testing"
)

func randString() string {
	res := ""
	for i := 0; i < 15+rand.Intn(20); i++ {
		res += string(rune(rand.Intn(100)))
	}
	return res
}

func TestHeap(t *testing.T) {
	heap := NewHeap(100)
	dic := make(map[string]int)
	for i := 0; i < 100; i++ {
		str := randString()
		_, ok := dic[str]
		if !ok {
			dic[str] = rand.Intn(100)
		}
	}

	for k := range dic {
		element := NewElement(k, "val")
		element.ref = uint8(rand.Intn(100))
		heap.Add(element)
	}

	if heap.verify(1) != true {
		log.Fatalf("Heap Verify Failed")
	}
}
