package lru

import (
	"fmt"
)

// LRU 双向链表LRU
type LRU struct {
	data map[string]*Node
	cap  int
	head *Node
}

func newLRU(cap int) *LRU {
	head := NewNode(nil, nil)
	head.prev = head
	head.next = head
	return &LRU{
		data: make(map[string]*Node),
		cap:  cap,
		head: head,
	}
}

func (l *LRU) Add(k, v string) {
	_, ok := l.data[k]
	if ok {
		// 如果数据已经存在，那就当成一次访问，移到队头
		l.Get(k)
		return
	}
	node := NewNode(k, v)
	if l.cap == 0 {
		last := l.head.GetTail()
		last.Remove()
		delete(l.data, last.key.(string))
		l.cap++
	}
	l.cap--
	l.data[k] = node
	node.AppendToFront(l.head)
}

func (l *LRU) Get(k string) (string, bool) {
	node, ok := l.data[k]
	if !ok {
		return "", false
	}
	node.Remove()
	node.AppendToFront(l.head)
	return node.value.(string), true
}

func (l *LRU) length() (length int) {
	for tmp := l.head.next; tmp != l.head; tmp = tmp.next {
		length++
	}
	return
}

func (l *LRU) Print() {
	print("   ")
	for tmp := l.head.next; tmp != l.head; tmp = tmp.next {
		fmt.Printf("%s-->", tmp.value)
	}
	fmt.Println()
}
