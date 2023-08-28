package bloomfilter

import "sync/atomic"

type LockFreeBitArray struct {
	array    []uint64
	arrayLen uint64
}

func NewLockBitArray(bitCount uint64) *LockFreeBitArray {
	byteSz := (bitCount + 63) >> 6
	return &LockFreeBitArray{
		array:    make([]uint64, byteSz),
		arrayLen: byteSz,
	}
}

func (l *LockFreeBitArray) set(pos uint64) {
	bytePos := pos >> 6 % l.arrayLen
	bitPos := pos & 0x3f
	for {
		element := atomic.LoadUint64(&l.array[bytePos])
		if element>>bitPos&1 == 1 {
			break
		}
		if atomic.CompareAndSwapUint64(&l.array[bytePos], element, element|(1<<bitPos)) {
			break
		}
	}
}

func (l *LockFreeBitArray) get(pos uint64) uint8 {
	bytePos := pos >> 6 % l.arrayLen
	bitPos := pos & 0x3f
	element := atomic.LoadUint64(&l.array[bytePos])
	return uint8(element >> bitPos & 1)
}
