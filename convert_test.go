package sino

import (
	"testing"
	"unsafe"
)

func TestBytesToU32s(t *testing.T) {
	bs := []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0, 0, 0, 0x03, 0, 0, 0, 0x04, 0, 0, 0}
	t.Log("bytes is", bs)
	ints := BytesToU32s(bs)
	if ints[0] != 1 || ints[1] != 2 || ints[2] != 3 || ints[3] != 4 {
		t.Log("ints is", ints)
	}
}

func TestStructToBytes(test *testing.T) {
	t := T{0x01, 0, 0, 0}
	bs := StructToBytes(&t)
	if len(bs) != int(unsafe.Sizeof(t)) {
		test.Error("bytes from struct is", bs)
	}
}
