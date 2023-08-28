package lru

import (
	"fmt"
)

type Element struct {
	key   string
	value string
	next  *Element
	prev  *Element
}

type List struct {
	head *Element
}

// LRU 双向链表LRU
type LRU struct {
	data map[string]*Element
	cap  int
	list List
}

func newLRU(cap int) *LRU {
	head := &Element{
		value: "",
	}
	head.next = head
	head.prev = head
	return &LRU{
		data: make(map[string]*Element),
		cap:  cap,
		list: List{head: head},
	}
}

func (l *LRU) length() (length uint32) {
	tmp := l.list.head
	for tmp.next != l.list.head {
		length++
		tmp = tmp.next
	}
	return
}

func (l *LRU) Print() {
	print("   ")
	for tmp := l.list.head.next; tmp != l.list.head; tmp = tmp.next {
		fmt.Printf("%s-->", tmp.value)
	}
	fmt.Println()
}

func (l *LRU) Add(k, v string) {
	element, ok := l.data[k]
	if !ok {
		element = &Element{
			key:   k,
			value: v,
		}
		if l.cap == 0 {
			last := l.list.head.prev
			//fmt.Printf("   Pass (%s,%s)\n", last.key, last.value)
			delete(l.data, last.key)
			last.prev.next = last.next
			last.next.prev = last.prev
		}
		if l.cap > 0 {
			l.cap--
		}
		l.data[k] = element
	}
	if element.next == element {
		return
	}
	if element.prev != nil {
		element.prev.next = element.next
		element.next.prev = element.prev
	}
	element.next = l.list.head.next
	element.prev = l.list.head
	element.next.prev = element
	element.prev.next = element
}

func (l *LRU) Get(k string) string {
	element, ok := l.data[k]
	if !ok {
		return ""
	}
	head := l.list.head
	element.prev.next = element.next
	element.next.prev = element.prev

	element.next = head.next
	element.prev = head
	element.prev.next = element
	element.next.prev = element
	return element.value
}
