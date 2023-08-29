package cms

import (
	"FastKV/cache/util"
	"sync/atomic"
	"unsafe"
)

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

type CMS struct {
	bitArray    *LockFree4BitArray
	hashFuncNum int
	size        uint32
	windowSz    uint32
}

type stringStruct struct {
	Data unsafe.Pointer
	Len  int
}

// NewCMS NewBloom
// exceptInsertions 是预计需要保存的元素个数
// fpp：期望假阳性概率
func NewCMS(exceptInsertions int) *CMS {
	// 一共4个hash函数，一个元素需要4个counter存储
	return &CMS{
		bitArray:    NewLockBitArray(exceptInsertions << 2),
		hashFuncNum: 4,
		size:        0,
		windowSz:    uint32(exceptInsertions) * 10,
	}
}

func (c *CMS) Frequency(key string) int {
	ss := (*stringStruct)(unsafe.Pointer(&key))

	hash64 := uint64(memhash(ss.Data, 0, uintptr(ss.Len)))
	hash1 := uint32(hash64)
	hash2 := uint32(hash64 >> 4)
	combinedHash := uint64(hash1)

	frequency := uint8(15)
	for i := 0; i < c.hashFuncNum; i++ {
		frequency = util.MinUint8(frequency, c.bitArray.get(combinedHash))
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

	if c.size++; c.size == c.windowSz {
		c.reset()
		atomic.StoreUint32(&c.size, 0)
	}
}

// Reset half all the value in cms
func (c *CMS) reset() {
	c.bitArray.reset()
}
