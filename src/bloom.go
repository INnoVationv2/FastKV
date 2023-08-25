package utils

import (
	"math"
	"unsafe"
)

type Bloom struct {
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

// NewBloom
// exceptInsertions 是预计需要保存的元素个数
// fpp：期望假阳性概率
func NewBloom(exceptInsertions int, fpp float64) *Bloom {
	// 求出位数组大小
	bitSz := uint64(math.Ceil(float64(-exceptInsertions) * math.Log(fpp) / (math.Log(2) * math.Log(2))))
	//求出hash函数的个数
	k := int(math.Max(1, math.Ceil(float64(bitSz)/float64(exceptInsertions)*math.Log(2))))

	return &Bloom{
		bitArray:    NewLockBitArray(bitSz),
		hashFuncNum: k,
	}
}

func (b *Bloom) Put(val string) {
	ss := (*stringStruct)(unsafe.Pointer(&val))

	hash64 := uint64(memhash(ss.Data, 0, uintptr(ss.Len)))
	hash1 := uint32(hash64)
	hash2 := uint32(hash64 >> 4)
	combinedHash := uint64(hash1)

	for i := 0; i < b.hashFuncNum; i++ {
		b.bitArray.set(combinedHash)
		combinedHash += uint64(hash2)
	}
}

func (b *Bloom) MayContain(val string) bool {
	ss := (*stringStruct)(unsafe.Pointer(&val))

	hash64 := uint64(memhash(ss.Data, 0, uintptr(ss.Len)))
	hash1 := uint32(hash64)
	hash2 := uint32(hash64 >> 4)
	combinedHash := uint64(hash1)

	for i := 0; i < b.hashFuncNum; i++ {
		if b.bitArray.get(combinedHash) == 0 {
			return false
		}
		combinedHash += uint64(hash2)
	}
	return true
}
