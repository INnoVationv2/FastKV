package corekv_diy

type Element struct {
	key   string
	value string
	next  *Element
	prev  *Element
}

type List struct {
	head *Element
}
