package cache

import (
	"math/rand"
)

func CeilingPowerOfTwo(i int) int {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}

func MinUint8(x, y uint8) uint8 {
	if x <= y {
		return x
	}
	return y
}

func RandString() string {
	res := ""
	for i := 0; i < 15+rand.Intn(20); i++ {
		res += string(rune(rand.Intn(100)))
	}
	return res
}

type Set map[string]*struct{}

func GetRandomStringSet(sz int) Set {
	set := make(Set)
	for i := 0; i < sz; i++ {
		str := RandString()
		_, ok := set[str]
		if ok {
			continue
		}
		set[str] = nil
	}
	return set
}
