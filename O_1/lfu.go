package O1

type LFU struct {
	head *Node
	dict map[string]*Element
	cap  uint32
}

func NewLFU(cap uint32) *LFU {
	node := &Node{}
	node.prev = node
	node.next = node
	return &LFU{
		head: node,
		dict: make(map[string]*Element),
		cap:  cap,
	}
}

func (l *LFU) Add(k, v string) {
	_, ok := l.dict[k]
	if ok {
		l.Get(k)
		return
	}

	if l.cap > 0 {
		l.cap--
	} else {
		node := l.head.next
		key := ""
		for key = range node.child {
			break
		}
		elem := node.child[key]
		node.Del(elem)

		delete(l.dict, elem.key)

		//fmt.Printf("Pass %s--%d\n", elem.key, elem.ref)
		elem.parent = nil
	}

	elem := NewElement(k, v)
	if l.head.next == l.head || l.head.next.ref != 0 {
		AddNode(l.head, NewNode(0))
	}

	l.dict[elem.key] = elem
	node := l.head.next
	node.Add(elem)
}

func AddNode(prevNode, node *Node) {
	node.next = prevNode.next
	node.prev = prevNode
	prevNode.next.prev = node
	prevNode.next = node
}

func (l *LFU) Get(k string) (string, bool) {
	element, ok := l.dict[k]
	if !ok {
		return "", false
	}

	element.ref++
	node := element.parent
	if node.next.ref != element.ref {
		AddNode(node, NewNode(element.ref))
	}
	newNode := node.next
	newNode.Add(element)

	node.Del(element)
	if node.childSz == 0 {
		node.prev.next = node.next
		node.next.prev = node.prev
	}

	return element.value, true
}
