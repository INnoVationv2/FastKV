package utils

import "math/rand"

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
