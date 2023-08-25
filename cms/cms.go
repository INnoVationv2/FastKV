package cms

import (
	"math"
	"unsafe"
)

type CMS struct {
	bitArray    *LockFreeBitArray
	hashFuncNum int
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

type stringStruct struct {
	Data unsafe.Pointer
	Len  int
}

// NewCMS NewBloom
// exceptInsertions 是预计需要保存的元素个数
// fpp：期望假阳性概率
func NewCMS(exceptInsertions int) *CMS {
	exceptInsertions <<= 2
	return &CMS{
		bitArray:    NewLockBitArray(exceptInsertions),
		hashFuncNum: 4,
	}
}

func (c *CMS) Frequency(key string) int {
	ss := (*stringStruct)(unsafe.Pointer(&key))

	hash64 := uint64(memhash(ss.Data, 0, uintptr(ss.Len)))
	hash1 := uint32(hash64)
	hash2 := uint32(hash64 >> 4)
	combinedHash := uint64(hash1)

	frequency := uint8(math.MaxUint8)
	for i := 0; i < c.hashFuncNum; i++ {
		frequency = minUint8(frequency, c.bitArray.get(combinedHash))
		combinedHash += uint64(hash2)
	}
	return int(frequency)
}

func (c *CMS) Increment(key string) {
	ss := (*stringStruct)(unsafe.Pointer(&key))

	hash64 := uint64(memhash(ss.Data, 0, uintptr(ss.Len)))
	hash1 := uint32(hash64)
	hash2 := uint32(hash64 >> 4)
	combinedHash := uint64(hash1)

	for i := 0; i < c.hashFuncNum; i++ {
		c.bitArray.incrementAt(combinedHash)
		combinedHash += uint64(hash2)
	}
}

// Reset half all the value in cms
func (c *CMS) Reset() {
	c.bitArray.reset()
}
