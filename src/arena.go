package utils

import (
	"github.com/pkg/errors"
	"log"
	"sync/atomic"
	"unsafe"
)

type Arena struct {
	// current used memory
	size     uint32
	capacity uint32
	buf      []byte
}

const (
	nodeAlign = int(unsafe.Sizeof(uint64(0))) - 1
)

func newArena(sz uint32) *Arena {
	return &Arena{
		size:     1,
		capacity: sz,
		buf:      make([]byte, sz),
	}
}

// Parameter: size which we want to allocate
// Return: allocated size start address
func (a *Arena) allocate(sz uint32) uint32 {
	offset := atomic.AddUint32(&a.size, sz)

	if a.capacity < offset {
		growSize := a.capacity
		if growSize > 1<<30 {
			growSize = 1 << 30
		}

		if growSize < sz {
			growSize = sz
		}

		newSize := atomic.AddUint32(&a.capacity, growSize)
		newBuf := make([]byte, newSize)
		AssertTrue(a.capacity == uint32(copy(newBuf, a.buf)))
		a.buf = newBuf
	}
	return offset - sz
}

// 传入要创建的Node的高度，返回存储该Node的起始地址
// Node默认以最大高度48创建，大部分节点没有那么高
// 因此可以以实际高度创建node中的tower，节省空间
func (a *Arena) putNode(height int) uint32 {
	nodeSz := int(unsafe.Sizeof(node{}))

	sz := int(unsafe.Sizeof(uint32(0)))
	unusedSize := sz * (defaultMaxLevel - height)

	realSize := uint32(nodeSz - unusedSize + nodeAlign)
	addr := a.allocate(realSize)
	return (addr + uint32(nodeAlign)) &^ uint32(nodeAlign)
}

func (a *Arena) putKey(key []byte) uint32 {
	keySz := uint32(len(key))
	offset := a.allocate(keySz)
	buf := a.buf[offset : offset+keySz]
	AssertTrue(len(key) == copy(buf, key))
	return offset
}

func (a *Arena) putVal(v ValueStruct) uint32 {
	sz := v.EncodedSize()
	offset := a.allocate(sz)
	v.EncodeValue(a.buf[offset:])
	return offset
}

// AssertTrue asserts that b is true. Otherwise, it would log fatal.
func AssertTrue(b bool) {
	if !b {
		log.Fatalf("%+v", errors.Errorf("Assert failed"))
	}
}
