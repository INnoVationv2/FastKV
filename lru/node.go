package lru

type Node struct {
	key   any
	value any
	next  *Node
	prev  *Node
}

func NewNode(key, val interface{}) *Node {
	node := &Node{}
	if key != nil {
		node.key = key
	}
	if val != nil {
		node.value = val
	}
	return node
}

func (n *Node) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (n *Node) GetTail() *Node {
	return n.prev
}

func (n *Node) AppendToFront(head *Node) {
	n.prev = head
	n.next = head.next
	n.prev.next = n
	n.next.prev = n
}
