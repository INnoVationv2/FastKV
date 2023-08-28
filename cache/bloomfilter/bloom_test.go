package bloomfilter

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func randString() string {
	res := ""
	for i := 0; i < 15+rand.Intn(20); i++ {
		res += string(rune(rand.Intn(100)))
	}
	return res
}

func TestBloom(t *testing.T) {
	testSz := int(1e6)
	set := make(map[string]int, testSz)
	for len(set) < testSz {
		set[randString()] = 0
	}

	bloom := NewBloom(testSz, 0.01)

	for k := range set {
		bloom.Put(k)
	}

	// 检验是否能够完成过滤
	for k := range set {
		if !bloom.MayContain(k) {
			log.Fatalf("Str Should exist!")
		}
	}

	// 测试误报率
	fp := 0
	for i := 0; i < testSz; {
		str := randString()
		_, ok := set[str]
		if !ok {
			if bloom.MayContain(str) {
				fp++
			}
			i++
		}
	}
	fmt.Printf("Fpp :%f\n", float64(fp)/float64(testSz))
}

func put(bloom *Bloom, val string) {
	bloom.Put(val)
}

func TestConcurrentBloom(t *testing.T) {
	testSz := int(1e6)
	set := make(map[string]int, testSz)
	for len(set) < testSz {
		set[randString()] = 0
	}

	bloom := NewBloom(testSz, 0.01)

	for k := range set {
		go put(bloom, k)
	}

	// 检验是否能够完成过滤
	for k := range set {
		if !bloom.MayContain(k) {
			log.Fatalf("Str Should exist!")
		}
	}

	// 测试误报率
	fp := 0
	for i := 0; i < testSz; {
		str := randString()
		_, ok := set[str]
		if !ok {
			if bloom.MayContain(str) {
				fp++
			}
			i++
		}
	}
	fmt.Printf("Fpp :%f\n", float64(fp)/float64(testSz))
}
