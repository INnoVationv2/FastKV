package utils

import (
	"bytes"
	"math/rand"
	"sync"
)

const (
	MaxHeight = 48
)

type Node struct {
	entry  *Entry
	next   []*Node
	height int
}

type SkipList struct {
	header *Node
	lock   sync.RWMutex
	height int
	size   int
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: &Node{
			entry:  nil,
			next:   make([]*Node, MaxHeight),
			height: MaxHeight,
		},
		height: MaxHeight,
		size:   0,
	}
}

func (entry *Entry) compare(entry2 *Entry) int {
	if entry.score < entry2.score {
		return -1
	} else if entry.score > entry2.score {
		return 1
	}

	return bytes.Compare(entry.key, entry2.key)
}

func compareKey(entry *Entry, key []byte) int {
	score := calcScore(key)
	if entry.score < score {
		return -1
	} else if entry.score > score {
		return 1
	}

	return bytes.Compare(entry.key, key)
}

func (list *SkipList) Add(newEntry *Entry) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	prevs := make([]*Node, MaxHeight)
	prev := list.header
	newEntry.score = calcScore(newEntry.key)

	for level := list.height - 1; level >= 0; level-- {
		for {
			if prev.next[level] == nil {
				break
			}

			if prev.next[level] == nil {
				break
			}
			cmp := prev.next[level].entry.compare(newEntry)
			if cmp == 0 {
				prev.next[level].entry.value = newEntry.value
				return true
			}
			if cmp > 0 {
				break
			}
			prev = prev.next[level]
		}
		prevs[level] = prev
	}

	newNode := NewNode(newEntry, randomHeight())
	for level := 0; level < newNode.height; level++ {
		newNode.next[level] = prevs[level].next[level]
		prevs[level].next[level] = newNode
	}
	return true
}

func (list *SkipList) Search(key []byte) *Entry {
	list.lock.RLock()
	defer list.lock.RUnlock()
	prev := list.header
	for level := list.height - 1; level >= 0; level-- {
		if prev.next[level] == nil {
			continue
		}

		for {
			if prev.next[level] == nil {
				break
			}
			cmp := compareKey(prev.next[level].entry, key)
			if cmp == 0 {
				return prev.next[level].entry
			}
			if cmp > 0 {
				break
			}
			prev = prev.next[level]
		}
	}
	return nil
}

func randomHeight() int {
	height := 1
	for i := 0; i < MaxHeight; i++ {
		if rand.Intn(2) == 1 {
			break
		}
		height++
	}
	return height
}

func NewNode(entry *Entry, height int) *Node {
	return &Node{
		entry:  entry,
		next:   make([]*Node, height),
		height: height,
	}
}

func calcScore(key []byte) uint64 {
	if key == nil {
		return 0
	}
	var hash uint64

	length := len(key)

	if length > 8 {
		length = 8
	}

	for i := 0; i < length; i++ {
		shift := 64 - 8*(i+1)
		hash |= uint64(key[i]) << shift
	}

	return hash
}
