package src

import (
	"github.com/pkg/errors"
	"log"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type Arena struct {
	height     uint32
	buffer     []*Buffer
	bufferLock sync.Mutex
}

type Buffer struct {
	sz       uint32
	capacity uint32
	buf      []byte
}

const (
	nodeAlign = int(unsafe.Sizeof(uint64(0))) - 1
)

func (a *Arena) getHeight() uint32 {
	return atomic.LoadUint32(&a.height)
}

func (a *Arena) getBuffer() *Buffer {
	return a.buffer[a.getHeight()-1]
}

func newArena(sz uint32) *Arena {
	arena := &Arena{
		height: 1,
		buffer: make([]*Buffer, 0),
	}
	arena.buffer = append(arena.buffer, newBuffer(sz))
	return arena
}

func newBuffer(sz uint32) *Buffer {
	return &Buffer{
		sz:       1,
		capacity: sz,
		buf:      make([]byte, sz),
	}
}

//func (a *Buffer) getAddr() uint64 {
//	return uint64(uintptr(unsafe.Pointer(a)))
//}

// Parameter: size which we want to allocate
// Return: allocated size start address
func (a *Arena) allocate(sz uint32) (uint32, uint32) {
	for {
		buf := a.getBuffer()
		height := a.height
		offset := atomic.AddUint32(&buf.sz, sz)
		capacity := buf.capacity
		if capacity > offset {
			return height - 1, offset - sz
		}

		for height == a.height {
			if TryLockWithTimeout(&a.bufferLock, time.Millisecond*10) {
				break
			}
		}

		if height != a.height {
			continue
		}

		growSize := capacity
		if growSize > 1<<30 {
			growSize = 1 << 30
		}

		if growSize < sz {
			growSize = sz
		}
		newSize := capacity + growSize
		newBuf := newBuffer(newSize)
		height = atomic.AddUint32(&a.height, 1)
		//a.buffer[height-1] = newBuf
		a.buffer = append(a.buffer, newBuf)

		a.bufferLock.Unlock()
	}
}

//
//func (a *Arena) setBuf(oldBuf uint64, newBuf *Buffer) bool {
//	return atomic.CompareAndSwapUint64(&a.bufAddr, oldBuf, newBuf.getAddr())
//}

//func getBuf(ptr uintptr) []byte {
//
//	return unsafe.Pointer(ptr)
//}

// 传入要创建的Node的高度，返回存储该Node的起始地址

func (a *Arena) putNode(h int) uint64 {
	nodeSz := int(unsafe.Sizeof(Node{}))

	// Node默认以最大高度48创建，大部分节点没有那么高
	// 因此可以以实际高度创建node中的tower，节省空间
	sz := int(unsafe.Sizeof(uint32(0)))
	unusedSize := sz * (maxHeight - h)
	realSize := uint32(nodeSz - unusedSize + nodeAlign)

	level, addr := a.allocate(realSize)
	addr = (addr + uint32(nodeAlign)) &^ uint32(nodeAlign)
	return uint64(level)<<56 | uint64(addr)
}

func (a *Arena) getNode(offset uint64) *Node {
	if offset == 0 {
		return nil
	}
	level := offset >> 56
	offset = offset << 8 >> 8
	return (*Node)(unsafe.Pointer(&a.buffer[level].buf[offset]))
}

func (a *Arena) putKey(key []byte) (uint32, uint32) {
	keySz := uint32(len(key))
	level, offset := a.allocate(keySz)
	buf := a.buffer[level].buf[offset : offset+keySz]
	AssertTrue(len(key) == copy(buf, key))
	return level, offset
}

func (a *Arena) getNodeOffset(nd *Node) uint64 {
	if nd == nil {
		return 0 //返回空指针
	}
	return nd.offset
}

func (a *Arena) getKey(offset uint32, size uint32) []byte {
	level := size >> 24
	size = size << 8 >> 8
	return a.buffer[level].buf[offset : offset+size]
}

func (a *Arena) putVal(v ValueStruct) (uint32, uint32) {
	sz := v.EncodedSize()
	level, offset := a.allocate(sz)
	v.EncodeValue(a.buffer[level].buf[offset:])
	return level, offset
}

func (a *Arena) getVal(valOffset uint32, valSize uint32) []byte {
	level := valSize >> 24
	valSize = valSize << 8 >> 8
	return a.buffer[level].buf[valOffset : valOffset+valSize]
}

// AssertTrue asserts that b is true. Otherwise, it would be fatal.
func AssertTrue(b bool) {
	if !b {
		log.Fatalf("%+v", errors.Errorf("Assert failed"))
	}
}
