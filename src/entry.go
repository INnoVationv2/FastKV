package FastKV

type Entry struct {
	key   []byte
	value []byte
	// For fast compare
	score uint64
}

func NewEntry(key []byte, value []byte) *Entry {
	return &Entry{
		key:   key,
		value: value,
		score: calcScore(key),
	}
}
