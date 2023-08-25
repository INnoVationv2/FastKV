package src

import (
	"sync/atomic"
	"unsafe"
)

const nodeSize = unsafe.Sizeof(Node{})

type Node struct {
	offset uint64
	//   value offset: uint32 (bits 0-31)
	//   value size  : uint16 (bits 32-63)
	value uint64

	// A byte slice is 24 bytes. We are trying to save space here.
	keyOffset uint32
	keySize   uint32

	// Height of the tower.
	height uint32

	tower [maxHeight]uint64
}

// 这里有一个问题，putNode时空间够用，node放在了旧的内存池
// 而putKey和PutVal时，发生了扩容，Key和Value分配在了新内存，
// 此时node也需要更新，指向新内存池
func newNode(arena *Arena, key []byte, value ValueStruct, height uint32) *Node {
	offset := arena.putNode(int(height))
	node := arena.getNode(offset)
	node.offset = offset
	node.height = height

	level, keyOffset := arena.putKey(key)
	node.keyOffset = keyOffset
	node.keySize = level<<24 | uint32(len(key))

	level, valOffset := arena.putVal(value)
	node.setValue(encodedValue(valOffset, level, value.EncodedSize()))
	return node
}

func encodedValue(valOffset uint32, level uint32, valSize uint32) uint64 {
	valSize = level<<24 | valSize
	return uint64(valSize)<<32 | uint64(valOffset)
}

func decodedValue(value uint64) (valOffset uint32, valSize uint32) {
	valOffset = uint32(value)
	valSize = uint32(value >> 32)
	return
}

func (n *Node) getValueOffset() (uint32, uint32) {
	return decodedValue(n.value)
}

func (n *Node) decodeValue() (offset uint32, size uint16) {
	value := atomic.LoadUint64(&n.value)
	offset = uint32(value)
	size = uint16(value >> 32)
	return
}

func (n *Node) getKey(arena *Arena) []byte {
	return arena.getKey(n.keyOffset, n.keySize)
}

func (n *Node) setValue(value uint64) {
	atomic.StoreUint64(&n.value, value)
}

func (n *Node) getNxtOffset(h int) uint64 {
	return atomic.LoadUint64(&n.tower[h])
}

func (n *Node) casNxtOffset(h int, old uint64, new uint64) bool {
	//log.Printf("[Node] casNxtOffset %d %d->%d\n", h, old, new)
	return atomic.CompareAndSwapUint64(&n.tower[h], old, new)
}

type ValueStruct struct {
	Value []byte
}

// EncodedSize 求出存储V需要多少空间
func (v *ValueStruct) EncodedSize() uint32 {
	return uint32(len(v.Value))
}

// EncodeValue 将v序列化到b中
func (v *ValueStruct) EncodeValue(buf []byte) uint32 {
	valueSz := copy(buf, v.Value)
	// 返回占用多少空间
	return uint32(valueSz)
}

// DecodeValue 从buf中解析ValueStruct
func (v *ValueStruct) DecodeValue(buf []byte) {
	v.Value = buf
}
