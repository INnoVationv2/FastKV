package lru

import "fmt"

type Set map[int]Pair
type Pair struct {
	key   string
	value string
}

func GetRandomString(sz int) Set {
	strSet := make(Set)
	for i := 0; i < sz; i++ {
		key := fmt.Sprintf("Key %d", i)
		value := fmt.Sprintf("Val %d", i)
		strSet[i] = Pair{key: key, value: value}
	}
	return strSet
}
