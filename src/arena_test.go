package src

import (
	"fmt"
	"testing"
	"unsafe"
)

func Test_1(t *testing.T) {
	fmt.Println("Checking uint32")
	size := unsafe.Sizeof(uint32(0))
	for i := 1; i < 10000; i++ {
		slice := make([]uint32, i)
		addr := uintptr(unsafe.Pointer(&slice[0]))
		fmt.Printf("%x\n", addr)
		if addr&size != 0 {
			fmt.Printf("Check Error: %v\n", addr)
		}
	}

	fmt.Println("Checking uint64")
	size = unsafe.Sizeof(uint64(0))
	for i := 1; i < 10000; i++ {
		slice := make([]uint64, i)
		addr := uintptr(unsafe.Pointer(&slice[0]))
		fmt.Printf("%x\n", addr)
		if addr&size != 0 {
			fmt.Printf("Check Error: %v\n", addr)
		}
	}
}

func Test_2(t *testing.T) {
	buf := make([]byte, 100)
	str := "hello, world"
	copy(buf, str)
	tmp := buf[len(str) : len(str)+len(str)]
	copy(tmp, str)

	buf2 := make([]byte, 200)
	res := copy(buf2, buf)
	println(res)
}
