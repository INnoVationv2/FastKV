package O_LogN

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestLFU(t *testing.T) {
	lfu := NewLFU(10)
	for i := 0; i < 100; i++ {
		k, v := "Key"+strconv.Itoa(i), "Val"+strconv.Itoa(i)
		lfu.Add(k, v)
		for j := 0; j < i; j++ {
			lfu.Get(k)
		}
		if i >= 10 {
			if len(lfu.dict) != 10 {
				log.Fatalf("dic Size not correct")
			}
		}
	}
}

// BenchmarkLFU_Add-10	1	2462610833 ns/op
func Test_LFU_Add(t *testing.T) {
	lfu := NewLFU(1e4)
	st := time.Now()
	for i := 0; i < 100000000; i++ {
		k, v := "Key"+strconv.Itoa(i), "Val"+strconv.Itoa(i)
		lfu.Add(k, v)
	}
	fmt.Println(time.Now().Sub(st))
}
