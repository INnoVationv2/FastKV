package O1

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestLFU(t *testing.T) {
	lfu := NewLFU(1e4)
	for i := 0; i < 10000000; i++ {
		k, v := "Key"+strconv.Itoa(i), "Val"+strconv.Itoa(i)
		lfu.Add(k, v)

		for j := 0; j < rand.Intn(500); j++ {
			lfu.Get("Key" + strconv.Itoa(rand.Intn(i+1)))
		}

		if i >= 10 {
			if len(lfu.dict) != 10 {
				log.Fatalf("dic Size not correct")
			}
		}
	}
}

// 24.552748958s
func Test_LFU_Add(t *testing.T) {
	lfu := NewLFU(1e4)
	st := time.Now()
	for i := 0; i < 100000000; i++ {
		k, v := "Key"+strconv.Itoa(i), "Val"+strconv.Itoa(i)
		lfu.Add(k, v)

		//for j := 0; j < rand.Intn(500); j++ {
		//	lfu.Get("Key" + strconv.Itoa(rand.Intn(i+1)))
		//}
	}
	fmt.Println(time.Now().Sub(st))
}
