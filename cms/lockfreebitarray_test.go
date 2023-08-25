package cms

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func Test_LockFreeArray(t *testing.T) {
	array := NewLockBitArray(1000)

	for i := 0; i < 1000; i++ {
		for j := 0; j < 10; j++ {
			array.incrementAt(uint64(i))
		}
	}

	for i := uint64(0); i < 1000; i++ {
		val := array.get(i)
		if val != uint8(10) {
			log.Fatalf("Not Equal: %d", i)
		}
	}
}

func inc(group *sync.WaitGroup, array *LockFreeBitArray, i uint64) {
	array.incrementAt(i)
	group.Done()
}

func Test_Concurrency_LFA(t *testing.T) {
	fail := 0
	for i := 0; i < 1000; i++ {
		array := NewLockBitArray(10000)
		insert(array)
		fail += array.fail
	}
	fmt.Printf("Average Fail: %d\n", fail/10000)
}

func insert(array *LockFreeBitArray) {
	group := &sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		for j := 0; j < 10; j++ {
			group.Add(1)
			go inc(group, array, uint64(i))
		}
	}
	group.Wait()

	for i := uint64(0); i < 10000; i++ {
		val := array.get(i)
		if val != uint8(10) {
			log.Fatalf("Not Equal: %d:%d", i, val)
		}
	}
}
