package sino

import (
	"reflect"
	"unsafe"
)

func BytesToU64s(b []byte) []uint64 {
	p := unsafe.Pointer(&b)
	pints := (*[]uint64)(p)
	(*pints) = (*pints)[0 : len(b)/8]
	return *pints
}

func BytesToU32s(b []byte) []uint32 {
	p := unsafe.Pointer(&b)
	pints := (*[]uint32)(p)
	(*pints) = (*pints)[0 : len(b)/4]
	return *pints
}

func BytesToU16s(b []byte) []uint16 {
	p := unsafe.Pointer(&b)
	pints := (*[]uint16)(p)
	(*pints) = (*pints)[0 : len(b)/2]
	return *pints
}

type T struct {
	bs    byte
	bss   byte
	bsss  byte
	bssss byte
}

func StructToBytes(t *T) []byte {
	sl := &reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(t)),
		Len:  int(unsafe.Sizeof(*t)),
		Cap:  int(unsafe.Sizeof(*t)),
	}
	b := *(*[]byte)(unsafe.Pointer(sl))
	return b
	//   return *(*[]byte)(
	//     unsafe.Pointer(
	//       &(reflect.SliceHeader{
	// 	Data:uintptr(unsafe.Pointer(t)),
	// 		   Len:15,
	// 		   Cap:15 }     )))
}
