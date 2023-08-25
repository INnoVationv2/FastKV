package O1

type Node struct {
	prev    *Node
	next    *Node
	childSz uint32
	ref     uint32
	child   map[string]*Element
}

func (n *Node) Add(elem *Element) {
	n.childSz++
	n.child[elem.key] = elem
	elem.parent = n
}

func (n *Node) Del(elem *Element) {
	n.childSz--
	delete(n.child, elem.key)
}

type Element struct {
	key    string
	value  string
	ref    uint32
	parent *Node
}

func NewElement(k, v string) *Element {
	return &Element{
		key:   k,
		value: v,
		ref:   0,
	}
}

func NewNode(ref uint32) *Node {
	return &Node{
		childSz: 0,
		ref:     ref,
		child:   make(map[string]*Element),
	}
}
