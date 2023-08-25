package utils

import (
	"bytes"
	"fmt"
	"log"
	"sync/atomic"
	"unsafe"
)

const (
	maxHeight = 48
)

type SkipList struct {
	height     uint32
	arena      *Arena
	headOffset uint32
	ref        int32
}

func NewSkipList() *SkipList {
	arena := newArena(uint32(nodeSize) * 1000)
	head := newNode(arena, nil, ValueStruct{}, maxHeight)
	return &SkipList{
		height:     1,
		arena:      arena,
		headOffset: head.offset,
		ref:        1,
	}
}

func (s *SkipList) getHead() *Node {
	return s.arena.getNode(s.headOffset)
}

func (s *SkipList) getHeight() uint32 {
	return atomic.LoadUint32(&s.height)
}

func (s *SkipList) findSpliceForLevel(key []byte, befOffset uint32, level int) (bef uint32, after uint32) {
	for {
		befNode := s.arena.getNode(befOffset)
		nxtOffset := atomic.LoadUint32(&befNode.tower[level])
		nxtNode := s.getNextNode(befNode, level)

		if nxtNode == nil {
			return befOffset, nxtOffset
		}

		nxtKey := nxtNode.getKey(s.arena)
		cmp := bytes.Compare(key, nxtKey)
		if cmp == 0 {
			return nxtOffset, nxtOffset
		}

		// prevKey < key < nxtKey
		if cmp < 0 {
			return befOffset, nxtOffset
		}

		befOffset = nxtOffset
	}
}

func (s *SkipList) findNear(key []byte, less bool, allowEqual bool) (target *Node, equal bool) {
	befNode := s.getHead()
	level := int(s.getHeight()) - 1
	for {
		nxtNode := s.getNextNode(befNode, level)
		if nxtNode == nil {
			// SkipList is Empty
			if befNode == s.getHead() {
				return nil, false
			}

			if level > 0 {
				level--
				continue
			}

			if !less {
				return nil, false
			}

			return befNode, false
		}

		nxtKey := s.arena.getKey(nxtNode.keyOffset, nxtNode.keySize)
		cmp := bytes.Compare(key, nxtKey)
		if cmp > 0 {
			befNode = nxtNode
			continue
		}

		if cmp == 0 {
			if allowEqual {
				return nxtNode, true
			}

			if !less {
				return s.getNextNode(nxtNode, 0), false
			}

			// 等于和大于情况都已处理，接下来是找node.key < key的情况
			if level > 0 {
				level--
				continue
			}

			return befNode, false
		}

		// cmp < 0
		// prevKey < key < nxtKey
		if level > 0 {
			level--
			continue
		}

		if !less {
			return nxtNode, false
		}

		return befNode, false
	}
}

func (s *SkipList) Search(key []byte) ValueStruct {
	node, equal := s.findNear(key, true, true) //找到小于等于key的节点
	if node == nil || !equal {
		return ValueStruct{}
	}
	node, equal = s.findNear(key, true, true)

	value := ValueStruct{}
	valOff, valSz := node.getValueOffset()
	value.Value = s.arena.getVal(valOff, valSz)
	if len(value.Value) == 0 {
		print("123")
	}
	return value
}

func (s *SkipList) getNextNode(n *Node, level int) *Node {
	nxtOffset := atomic.LoadUint32(&n.tower[level])
	return s.arena.getNode(nxtOffset)
}

func (s *SkipList) Add(e *Entry) {
	key, v := e.Key, ValueStruct{
		Value: e.Value,
	}

	listHeight := s.getHeight()
	var prev [maxHeight + 1]uint32
	var nxt [maxHeight + 1]uint32

	prev[listHeight] = s.headOffset
	for level := int(listHeight) - 1; level >= 0; level-- {
		prev[level], nxt[level] = s.findSpliceForLevel(key, prev[level+1], level)

		if prev[level] == nxt[level] {
			valOffset := s.arena.putVal(v)
			encValue := encodedValue(valOffset, v.EncodedSize())
			prevNode := s.arena.getNode(prev[level])
			prevNode.setValue(encValue)
			return
		}
	}

	height := randomHeight()
	node := newNode(s.arena, key, v, height)
	for {
		listHeight = s.getHeight()
		if height < listHeight || atomic.CompareAndSwapUint32(&s.height, listHeight, height) {
			break
		}
	}

	nodeOffset := s.arena.getNodeOffset(node)
	for level := uint32(0); level < height; level++ {
		for {
			prevNode := s.arena.getNode(prev[level])
			if prevNode == nil {
				prev[level], nxt[level] = s.findSpliceForLevel(key, s.headOffset, int(level))
				AssertTrue(prev[level] != nxt[level])
				continue
			}

			node.tower[level] = nxt[level]
			if prevNode.casNxtOffset(int(level), nxt[level], nodeOffset) {
				break
			}

			// cas失败，要重新寻找插入位置
			prev[level], nxt[level] = s.findSpliceForLevel(key, prev[level], int(level))
			// 当前值已经被别的插入了
			if prev[level] == nxt[level] {
				if level != 0 {
					log.Fatalf("Equality can happen only on base level: %d", level)
				}
				valOffset := s.arena.putVal(v)
				encValue := encodedValue(valOffset, v.EncodedSize())
				prevNode := s.arena.getNode(prev[level])
				prevNode.setValue(encValue)
				return
			}
		}
	}
}

func dynamicPrint(n int, ch rune) {
	for i := 0; i < n; i++ {
		fmt.Printf("%c", ch)
	}
}

func (s *SkipList) verify() bool {
	head := s.arena.getNode(s.headOffset)
	for node := s.arena.getNode(head.tower[0]); node != nil; {
		if len(node.getKey(s.arena)) == 0 || len(s.arena.getVal(node.getValueOffset())) == 0 {
			return false
		}
		node = s.arena.getNode(node.tower[0])
	}
	return true
}

func (s *SkipList) Draw() {
	baseNodes := make([]*Node, 0)
	head := s.arena.getNode(s.headOffset)
	for node := s.arena.getNode(head.tower[0]); node != nil; {
		baseNodes = append(baseNodes, node)
		node = s.arena.getNode(node.tower[0])
	}

	fmt.Printf("<")
	dynamicPrint(len(baseNodes)/2*20, '=')
	fmt.Printf("SkipList")
	dynamicPrint(len(baseNodes)/2*20, '=')
	fmt.Printf("==>\n")

	for level := int(s.getHeight()) - 1; level >= 0; level-- {
		prevNode := head
		fmt.Printf("   [%d]head", level)
		idx := 0
		for {
			node := s.arena.getNode(prevNode.tower[level])

			if node == nil {
				fmt.Printf("\n")
				break
			}

			for ; idx < len(baseNodes) && baseNodes[idx] != node; idx++ {
				dynamicPrint(20, '-')
			}

			nodePtr := uintptr(unsafe.Pointer(node))
			nxtNodePtr := uintptr(unsafe.Pointer(s.getNextNode(node, level)))
			fmt.Printf("-->")
			fmt.Printf("(%#05x,%#05x)", nodePtr&0xFFFFF, nxtNodePtr&0xFFFFF)
			//fmt.Printf("(%#05x,%#05x)", s.arena.getNodeOffset(node), node.tower[level])
			prevNode = node
			idx++
		}
	}

	fmt.Println()
	for i := 0; i < len(baseNodes); i++ {
		node := baseNodes[i]
		fmt.Printf("   (%#05x):{Height:%d, Key:%s, Value:%#v}\n",
			uintptr(unsafe.Pointer(node))&0xFFFFF,
			node.height,
			node.getKey(s.arena),
			s.arena.getVal(node.getValueOffset()))
	}

	fmt.Printf("<")
	dynamicPrint(len(baseNodes)/2*20, '=')
	fmt.Printf("===End===")
	dynamicPrint(len(baseNodes)/2*20, '=')
	fmt.Printf("==>\n")
}

type SkipListIterator struct {
	skipList  *SkipList
	curOffset uint32
}

func (s *SkipList) NewSkipListIterator() *SkipListIterator {
	return &SkipListIterator{
		skipList:  s,
		curOffset: s.getNextNode(s.getHead(), 0).offset,
	}
}

func (i *SkipListIterator) Rewind() {
	i.curOffset = i.skipList.getNextNode(i.skipList.getHead(), 0).offset
}

func (i *SkipListIterator) Valid() bool {
	return i.skipList.arena.getNode(i.curOffset) != nil
}

func (i *SkipListIterator) Next() {
	curNode := i.skipList.arena.getNode(i.curOffset)
	i.curOffset = curNode.tower[0]
}

func (i *SkipListIterator) Item() *Entry {
	node := i.skipList.arena.getNode(i.curOffset)
	key := i.skipList.arena.getKey(node.keyOffset, node.keySize)
	val := i.skipList.arena.getVal(decodedValue(node.value))
	return NewEntry(key, val)
}
