package FastKV

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sync"
	"testing"
)

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestSkipList_compare(t *testing.T) {
	key1 := []byte("1")
	key2 := []byte("2")
	node1 := NewEntry(key1, nil)
	node2 := NewEntry(key2, nil)
	assert.Equal(t, node1.compare(node1), 0)
	assert.Equal(t, node1.compare(node2), -1)
	assert.Equal(t, node2.compare(node1), 1)
}

func TestSkipListBasicCRUD(t *testing.T) {
	list := NewSkipList()
	// PUT & GET
	entry1 := NewEntry([]byte("Key1"), []byte("Val1"))
	assert.True(t, list.Add(entry1))
	assert.Equal(t, entry1.value, list.Search(entry1.key).value)

	entry2 := NewEntry([]byte("Key2"), []byte("Val2"))
	assert.True(t, list.Add(entry2))
	assert.Equal(t, entry2.value, list.Search(entry2.key).value)

	// Get a not exist entry
	assert.Nil(t, list.Search([]byte("noexist")))

	//Update an entry
	entry2New := NewEntry([]byte("Key1"), []byte("Val1+1"))
	assert.True(t, list.Add(entry2New))
	assert.Equal(t, entry2New.value, list.Search(entry2New.key).value)
}

func Benchmark_SkipListBasicCRUD(b *testing.B) {
	list := NewSkipList()
	key, val := "", ""
	maxTime := 1000000
	for i := 0; i < maxTime; i++ {
		key, val = fmt.Sprintf("Key%d", i), fmt.Sprintf("Val%d", i)
		entry := NewEntry([]byte(key), []byte(val))
		res := list.Add(entry)
		assert.True(b, res)
		searchVal := list.Search([]byte(key))
		assert.Equal(b, searchVal.value, []byte(val))
	}
}

func TestConcurrentBasic(t *testing.T) {
	const n = 1000
	l := NewSkipList()
	var wg sync.WaitGroup
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.True(t, l.Add(NewEntry(key(i), key(i))))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			if v != nil {
				require.EqualValues(t, key(i), v.value)
				return
			}
			require.Nil(t, v)
		}(i)
	}
	wg.Wait()
}

func Benchmark_ConcurrentBasic(b *testing.B) {
	const n = 1000
	l := NewSkipList()
	var wg sync.WaitGroup
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.True(b, l.Add(NewEntry(key(i), key(i))))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			if v != nil {
				require.EqualValues(b, key(i), v.value)
				return
			}
			require.Nil(b, v)
		}(i)
	}
	wg.Wait()
}
