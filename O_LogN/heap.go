package O_LogN

import "fmt"

// Heap 小根堆
type Heap struct {
	array []*Element
	// sz: 下一个元素要插入的位置，从1开始
	sz  uint32
	cap uint32
}

func NewHeap(cap int) *Heap {
	return &Heap{
		array: make([]*Element, cap+1),
		sz:    1,
		cap:   uint32(cap + 1),
	}
}

func (h *Heap) Add(element *Element) {
	element.pos = h.sz
	h.array[h.sz] = element
	h.up(h.sz)
	h.sz++
}

func (h *Heap) Del() {
	h.swap(1, h.sz-1)
	h.sz--
	h.down(1)
}

func (h *Heap) up(pos uint32) {
	if pos < 1 || pos >= h.sz {
		return
	}
	element := h.array[pos]
	for {
		upPos := pos >> 1
		if upPos < 1 {
			break
		}

		upElement := h.array[upPos]
		if upElement.ref <= element.ref {
			break
		}

		h.swap(pos, upPos)
		pos = upPos
	}
}

func (h *Heap) down(pos uint32) {
	if pos < 1 || pos >= h.sz {
		return
	}
	for {
		element := h.array[pos]
		targetElem := element

		if lChild := pos << 1; lChild < h.sz &&
			h.array[lChild].ref < targetElem.ref {
			targetElem = h.array[lChild]
		}

		if rChild := pos<<1 + 1; rChild < h.sz &&
			h.array[rChild].ref < targetElem.ref {
			targetElem = h.array[rChild]
		}
		if targetElem == element {
			break
		}
		pos = targetElem.pos
		h.swap(element.pos, targetElem.pos)
	}
}

func (h *Heap) swap(pos1, pos2 uint32) {
	h.array[pos1].pos, h.array[pos2].pos = h.array[pos2].pos, h.array[pos1].pos
	h.array[pos1], h.array[pos2] = h.array[pos2], h.array[pos1]
}

func (h *Heap) growSz() {
	newCap := h.cap << 1
	newArray := make([]*Element, newCap)
	copy(newArray, h.array)
	h.array = newArray
	h.cap = newCap
}

// 检验是否为小根堆
func (h *Heap) verify(idx uint32) bool {
	if idx >= h.sz {
		return true
	}

	ref := h.array[idx].ref
	res := true

	if lChild := idx << 1; lChild < h.sz {
		if h.array[lChild].ref >= ref {
			res = h.verify(lChild)
		} else {
			return false
		}
	}

	if rChild := idx<<1 + 1; res && rChild < h.sz {
		if h.array[rChild].ref >= ref {
			res = h.verify(rChild)
		} else {
			res = false
		}
	}
	return res
}

func (h *Heap) print() {
	length := uint32(1)
	fmt.Println("=======================")
	for i := uint32(1); i < h.sz; length <<= 1 {
		for j := i + length; i < h.sz && i < j; i++ {
			fmt.Printf("%2d  ", h.array[i].ref)
		}
		fmt.Println()
	}
	fmt.Println("=======================")
}
