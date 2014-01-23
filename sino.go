package sino

import (
	"fmt"
)

func HexDump(addr uint64, data []byte) {
	length := len(data)
	nlines := length / 16
	if (length % 16) != 0 {
		nlines++
	}
	nlines--
	for i := 0; i < nlines; i++ {
		fmt.Printf("%X: ", addr)
		for j := 0; j < 16; j++ {
			fmt.Printf("%x ", data[i*16+j])
		}
		addr += 16
		fmt.Println("")
	}
	fmt.Printf("%X: ", addr)
	for j := 0; j < length%16; j++ {
		fmt.Printf("%x ", data[nlines*16+j])
	}
	fmt.Println("")
}
