package main

//package sino
import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"syscall"
)

const help_msg = `devmem Dump memory or memory mapped registers.
	Usage: devmem address(hexadecimal) [width(8, 16, 32)] value(hexadecimal)
			address: hexadecimal representation of the physical address(no '0x' prefix needed)
			width:   8, 16, 32 bits
			value:   hexadecimal representation of the value to be written(no '0x' prefix needed)
`

//func devmem(addr int64, width uint8, value uint32) {
func main() {
	var addr int64
	var err error
	var width uint64=32
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
		log.Printf("value is %x\n",value)
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
			log.Fatal("Error open /dev/mem.\n"+err.Error())
		}
		defer syscall.Close(fd)
		//func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
		log.Printf("mmap start addr is %x\n", addr&(^0xFFF))
		map_data, err = syscall.Mmap(fd, addr&(^0xFFF), 0x1000, syscall.PROT_READ | syscall.PROT_WRITE, syscall.MAP_SHARED)
		if err != nil {
			log.Fatal("Error mmap /dev/mem.\n"+err.Error())
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
		case width == 32:
			map_data[index+3] = uint8(value >> 24)
			map_data[index+2] = uint8(value >> 16)
			fallthrough
		case width == 16:
			log.Printf("write %x",uint8(value>>8))
			map_data[index+1] = uint8(value >> 8)
			fallthrough
		case width == 8:
			log.Printf("write %x",uint8(value))
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
			fmt.Printf("%02x%02x \n", map_data[index+1], map_data[index])
		case width == 32:
			fmt.Printf("%02x%02x%02x%02x \n", map_data[index+3], map_data[index+2], map_data[index+1], map_data[index])
		default:
			log.Fatal("Wrong width.\n")
			fmt.Printf(help_msg)
		}
	}

}
