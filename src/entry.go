package src

// Entry _ 最外层写入的结构体
type Entry struct {
	Key   []byte
	Value []byte
}

// NewEntry_
func NewEntry(key, value []byte) *Entry {
	return &Entry{
		Key:   key,
		Value: value,
	}
}
