package utils

import "encoding/binary"

const defaultMaxLevel = 48

type node struct {
	// Multiple parts of the value are encoded as a single uint64 so that it
	// can be atomically loaded and stored:
	//   value offset: uint32 (bits 0-31)
	//   value size  : uint16 (bits 32-63)
	value uint64

	// A byte slice is 24 bytes. We are trying to save space here.
	keyOffset uint32
	keySize   uint16

	// Height of the tower.
	height uint16

	tower [defaultMaxLevel]uint32
}

func encodedValue(valOffset uint32, valSize uint32) uint64 {
	return uint64(valSize)<<32 | uint64(valOffset)
}

func decodedValue(value uint64) (valOffset uint32, valSize uint32) {
	valOffset = uint32(value)
	valSize = uint32(value >> 32)
	return
}

type ValueStruct struct {
	Value    []byte
	ExpireAt uint64
}

// EncodedSize 求出存储V需要多少空间
func (v *ValueStruct) EncodedSize() uint32 {
	valueSz := len(v.Value)
	expireAtSz := sizeVarint(v.ExpireAt)
	return uint32(valueSz + expireAtSz)
}

// EncodeValue 将v序列化到b中
func (v *ValueStruct) EncodeValue(buf []byte) uint32 {
	// 先放入ExpireAt
	expireAtSz := binary.PutUvarint(buf[:], v.ExpireAt)
	// 再放入value
	valueSz := copy(buf[expireAtSz:], v.Value)
	// 返回占用多少空间
	return uint32(expireAtSz + valueSz)
}

// DecodeValue 从buf中解析ValueStruct
func (v *ValueStruct) DecodeValue(buf []byte) {
	var size int
	v.ExpireAt, size = binary.Uvarint(buf)
	v.Value = buf[size:]
}

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
