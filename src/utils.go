package src

import (
	"math/rand"
	"sync"
	"time"
)

// 求出存储x需要多少Byte
func sizeVarint(x uint64) (sz int) {
	for {
		sz++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return
}

func randomHeight() uint32 {
	for i := 1; i <= maxHeight; i++ {
		if rand.Intn(2) == 0 {
			return uint32(i)
		}
	}
	return maxHeight
}

func TryLockWithTimeout(lock *sync.Mutex, timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		return false
	default:
		lock.Lock()
		return true
	}
}
