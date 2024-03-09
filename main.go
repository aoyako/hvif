package main

import (
	"encoding/binary"
	"fmt"
	"hvif/style"
	"io"
	"os"
	"unsafe"
)

func main() {
	filename := "res/folder.hvif"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	magic := make([]byte, 4)
	nread, err := io.ReadFull(file, magic)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	fmt.Printf("Magic %d bytes: %s\n", nread, magic)

	var count byte
	binary.Read(file, binary.LittleEndian, &count)
	fmt.Printf("Color count %d bytes: %d\n", unsafe.Sizeof(count), count)

	var i byte
	for i = 0; i < count; i++ {
		data := style.Read(file)
		if c, ok := data.Color(); ok {
			fmt.Println("color", c)
		}
		if g, ok := data.Gradient(); ok {
			fmt.Println("gradient", g)
		}
	}
}
