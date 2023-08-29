package lru

import (
	"fmt"
	"sync"
)

const (
	Window = iota
	Probation
	Protected
)

type Node struct {
	Type  int
	Key   string
	Value string
	next  *Node
	prev  *Node
	lock  *sync.RWMutex
}

func NewNode(key, val string) *Node {
	return &Node{
		Type:  Window,
		Key:   key,
		Value: val,
		lock:  nil,
	}
}

func NewEmptyNode() *Node {
	node := &Node{}
	node.prev = node
	node.next = node
	node.lock = &sync.RWMutex{}
	return node
}

func (n *Node) AppendToFront(head *Node) {
	n.prev = head
	n.next = head.next
	n.prev.next = n
	n.next.prev = n
}

func (n *Node) MoveToFront(head *Node) {
	n.Remove()
	n.AppendToFront(head)
}

func (n *Node) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (n *Node) GetTail() *Node {
	node := n.prev
	return node
}

func (n *Node) String() string {
	var str string
	for cur := n.next; cur != n; cur = cur.next {
		str += fmt.Sprintf("%s, ", cur.Value)
	}
	return str
}
