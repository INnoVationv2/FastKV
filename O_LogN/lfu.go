package O_LogN

type Element struct {
	key string
	val string
	ref uint8
	pos uint32
}

func NewElement(k, v string) *Element {
	return &Element{
		key: k,
		val: v,
		ref: 1,
	}
}

type LFU struct {
	cap  int
	dict map[string]*Element
	heap *Heap
}

func NewLFU(cap int) *LFU {
	return &LFU{
		cap:  cap,
		dict: make(map[string]*Element),
		heap: NewHeap(cap),
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
		//fmt.Printf("Pass %v--%v\n", l.heap.array[1].key, l.heap.array[1].ref)
		delete(l.dict, l.heap.array[1].key)
		l.heap.Del()
	}

	element := NewElement(k, v)
	l.dict[k] = element
	l.heap.Add(element)
}

func (l *LFU) Get(k string) (string, bool) {
	element, ok := l.dict[k]
	if !ok {
		return "", false
	}
	element.ref++
	l.heap.down(element.pos)
	return element.val, true
}
