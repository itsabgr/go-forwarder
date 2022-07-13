package handy

import (
	"golang.org/x/exp/constraints"
	"unsafe"
)

//MaxUint is max uint value
const MaxUint = ^uint(0)

//MaxInt is max int value
const MaxInt = int(MaxUint >> 1)

//UintptrSize is memory size of uintptr
const UintptrSize = unsafe.Sizeof(uintptr(0))

//Bits return bit size of any integer
func Bits[T constraints.Integer](n T) int {
	n = 0
	n = ^n
	i := 0
	for ; n<<i != 0; i++ {
	}
	return i
}
