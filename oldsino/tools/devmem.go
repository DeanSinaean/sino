package main

//package sino
import (
	"flag"
	"fmt"
	"log"
	"sino"
	"strconv"
	"syscall"
)

const help_msg = `devmem Dump memory or memory mapped registers.
	Usage: devmem address(hexadecimal) [width(8, 16, 32, 64)] value(hexadecimal)
			address: hexadecimal representation of the physical address(no '0x' prefix needed)
			width:   8, 16, 32 64 bits
			value:   hexadecimal representation of the value to be written(no '0x' prefix needed)
`

//func devmem(addr int64, width uint8, value uint32) {
func main() {
	var addr int64
	var err error
	var width uint64 = 32
	var value uint64
	var map_data []byte
	var index int
	flag.Parse()
	log.SetFlags(log.Lshortfile)
	log.Printf("NArg is %d\n", flag.NArg())
	switch {
	case flag.NArg() > 3:
		fmt.Printf("Too many arguments.\n")
		log.Fatal(help_msg)
	case flag.NArg() == 3:
		value, err = strconv.ParseUint(flag.Arg(2), 16, 32)
		if err != nil {
			fmt.Printf("Invalid arguments.\n")
			fmt.Printf(help_msg)
			log.Fatal("Invalid arguments.\n")
		}
		log.Printf("value is %x\n", value)
		fallthrough
	case flag.NArg() >= 2:
		width, err = strconv.ParseUint(flag.Arg(1), 10, 8)
		log.Printf("width is %d\n", width)
		if err != nil {
			fmt.Printf("Invalid arguments.\n")
			fmt.Printf(help_msg)
			log.Fatal("Invalid arguments.\n")
		}
		fallthrough
	case flag.NArg() >= 1:
		addr, err = strconv.ParseInt(flag.Arg(0), 16, 64)
		if err != nil {
			fmt.Printf("Invalid arguments.\n")
			fmt.Printf(help_msg)
			log.Fatal("Invalid arguments.\n")
		}
		log.Printf("addr is %x\n", addr)
		fd, err := syscall.Open("/dev/mem", syscall.O_RDWR|syscall.O_SYNC, 0)
		if err != nil {
			log.Fatal("Error open /dev/mem.\n" + err.Error())
		}
		defer syscall.Close(fd)
		//func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
		log.Printf("mmap start addr is %x\n", addr&(^0xFFF))
		map_data, err = syscall.Mmap(fd, addr&(^0xFFF), 0x1000, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			log.Fatal("Error mmap /dev/mem.\n" + err.Error())
		}
		defer syscall.Munmap(map_data)
		log.Printf("width is %d\n", width)
		index = int(addr & 0xFFF)
	default:
		fmt.Printf(help_msg)
	}
	log.Printf("width is %d\n", width)

	if flag.NArg() == 3 {
		switch {
		case width == 64:
			log.Printf("write %x", uint64(value))
			u64data := sino.BytesToU64s(map_data)
			u64data[index/8] = uint64(value)
		case width == 32:
			log.Printf("write %x", uint32(value))
			u32data := sino.BytesToU32s(map_data)
			u32data[index/4] = uint32(value)
		case width == 16:
			log.Printf("write %x", uint16(value))
			u16data := sino.BytesToU16s(map_data)
			u16data[index/2] = uint16(value)
		case width == 8:
			log.Printf("write %x", uint8(value))
			map_data[index] = uint8(value)
		default:
			log.Fatal("Wrong width.\n")
			fmt.Printf(help_msg)

		}

	} else {
		switch {
		case width == 8:
			fmt.Printf("%02x \n", map_data[index])
		case width == 16:
			u16data := sino.BytesToU16s(map_data)
			fmt.Printf("%04x \n", u16data[index/2])
		case width == 32:
			u32data := sino.BytesToU32s(map_data)
			fmt.Printf("%08x \n", u32data[index/4])
		case width == 64:
			u64data := sino.BytesToU64s(map_data)
			fmt.Printf("%016x \n", u64data[index/8])
		default:
			log.Fatal("Wrong width.\n")
			fmt.Printf(help_msg)
		}
	}

}
