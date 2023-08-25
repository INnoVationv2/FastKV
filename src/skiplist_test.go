// Copyright 2021 hardcore-os Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License")
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := rand.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestSkipList_Draw(t *testing.T) {
	list := NewSkipList()
	head := list.getHead()

	node1 := newNode(list.arena, []byte("key1"), ValueStruct{}, 5)
	node2 := newNode(list.arena, []byte("key2"), ValueStruct{}, 1)
	node3 := newNode(list.arena, []byte("key3"), ValueStruct{}, 3)
	node4 := newNode(list.arena, []byte("key4"), ValueStruct{}, 5)
	list.height = 5
	head.tower[4] = list.arena.getNodeOffset(node1)
	head.tower[3] = list.arena.getNodeOffset(node1)
	head.tower[2] = list.arena.getNodeOffset(node1)
	head.tower[1] = list.arena.getNodeOffset(node1)
	head.tower[0] = list.arena.getNodeOffset(node1)

	node1.tower[4] = list.arena.getNodeOffset(node4)
	node1.tower[3] = list.arena.getNodeOffset(node4)
	node1.tower[2] = list.arena.getNodeOffset(node3)
	node1.tower[1] = list.arena.getNodeOffset(node3)
	node1.tower[0] = list.arena.getNodeOffset(node2)

	node2.tower[0] = list.arena.getNodeOffset(node3)

	node3.tower[2] = list.arena.getNodeOffset(node4)
	node3.tower[1] = list.arena.getNodeOffset(node4)
	node3.tower[0] = list.arena.getNodeOffset(node4)

	list.Draw()
}

func TestSkipListBasicCRUD(t *testing.T) {
	list := NewSkipList()
	for i := 1; i < 100000; i++ {
		entry1 := NewEntry([]byte(RandString(10)), []byte("Val"+strconv.Itoa(i)))
		list.Add(entry1)
		assert.Equal(t, entry1.Value, list.Search(entry1.Key).Value)

		entry1.Value = []byte("New Val" + strconv.Itoa(i))
		list.Add(entry1)
		assert.Equal(t, entry1.Value, list.Search(entry1.Key).Value)
	}
}

func Benchmark_SkipListBasicCRUD(b *testing.B) {
	list := NewSkipList()
	for i := 1; i < 1000; i++ {
		entry1 := NewEntry([]byte(RandString(10)), []byte("Val"+strconv.Itoa(i)))
		list.Add(entry1)
		assert.Equal(b, entry1.Value, list.Search(entry1.Key).Value)

		entry1.Value = []byte("New Val" + strconv.Itoa(i))
		list.Add(entry1)
		assert.Equal(b, entry1.Value, list.Search(entry1.Key).Value)
	}
}

func TestConcurrentBasic(t *testing.T) {
	const n = 1000
	l := NewSkipList()
	var wg sync.WaitGroup
	key := func(i int) []byte {
		return []byte(fmt.Sprintf("Keykeykey%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Add(NewEntry(key(i), key(i)))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			require.EqualValues(t, key(i), v.Value)
			return

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
		return []byte(fmt.Sprintf("keykeykey%05d", i))
	}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Add(NewEntry(key(i), key(i)))
		}(i)
	}
	wg.Wait()

	// Check values. Concurrent reads.
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := l.Search(key(i))
			require.EqualValues(b, key(i), v.Value)
			require.NotNil(b, v)
		}(i)
	}
	wg.Wait()
}

func TestSkipListIterator(t *testing.T) {
	list := NewSkipList()

	//Put & Get
	entry1 := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry1)
	assert.Equal(t, entry1.Value, list.Search(entry1.Key).Value)

	entry2 := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry2)
	assert.Equal(t, entry2.Value, list.Search(entry2.Key).Value)

	//Update an entry
	entry2_new := NewEntry([]byte(RandString(10)), []byte(RandString(10)))
	list.Add(entry2_new)
	assert.Equal(t, entry2_new.Value, list.Search(entry2_new.Key).Value)

	iter := list.NewSkipListIterator()
	for iter.Rewind(); iter.Valid(); iter.Next() {
		fmt.Printf("iter key %s, value %s", iter.Item().Key, iter.Item().Value)
	}
	assert.Equal(t, 1, 1)
}
