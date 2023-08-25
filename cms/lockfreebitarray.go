package cms

import (
	"sync/atomic"
	"unsafe"
)

type LockFreeBitArray struct {
	array     []unsafe.Pointer
	arrayLen  int
	blockMask uint64

	fail int
}

func NewLockBitArray(entryNum int) *LockFreeBitArray {
	// 比entryNum大的2的幂
	sz := ceilingPowerOfTwo((entryNum + 1) >> 1)
	// 2的幂，必然是只有一个1，其余位都是0，减去1二进制就全是1，加速查找
	blockMask := sz - 1

	array := make([]unsafe.Pointer, sz)
	for i := 0; i < sz; i++ {
		var val uint8 = 0
		array[i] = unsafe.Pointer(&val)
	}

	return &LockFreeBitArray{
		array:     array,
		arrayLen:  sz,
		blockMask: uint64(blockMask),

		fail: 0,
	}
}

// 8bit, 每个元素占据4bit，每个能存储2个数据
func (l *LockFreeBitArray) incrementAt(pos uint64) {
	tableIdx := pos >> 1 & l.blockMask
	counterIdx := (pos & 0x01) << 2
	mask := uint8(0xf) << counterIdx

	for {
		oldCounterPtr := atomic.LoadPointer(&l.array[tableIdx])
		oldCounter := *(*uint8)(oldCounterPtr)

		if oldCounter&mask == mask {
			break
		}

		newCounter := oldCounter + (1 << counterIdx)
		if atomic.CompareAndSwapPointer(&l.array[tableIdx], oldCounterPtr, unsafe.Pointer(&newCounter)) {
			break
		}

		l.fail++
	}
}

func (l *LockFreeBitArray) get(pos uint64) uint8 {
	tableIdx := pos >> 1 & l.blockMask
	counterIdx := (pos & 0x01) << 2
	counter := *(*uint8)(atomic.LoadPointer(&l.array[tableIdx]))
	return (counter >> counterIdx) & 0xf
}

func (l *LockFreeBitArray) reset() {
	for i := 0; i < l.arrayLen; i++ {
		for {
			oldCounterPtr := atomic.LoadPointer(&l.array[i])
			oldVal := *(*uint8)(oldCounterPtr)
			if oldVal == 0 {
				break
			}

			newVal := oldVal >> 1 & 0x77
			if atomic.CompareAndSwapPointer(&l.array[i], oldCounterPtr, unsafe.Pointer(&newVal)) {
				break
			}
		}
	}
}
