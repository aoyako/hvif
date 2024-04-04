package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"

	"hvif/shape"
)

func main() {
	filename := "res/ime"

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
		s, err := readStyle(file)
		fmt.Printf("%+v, %v\n\n", s, err)
	}

	binary.Read(file, binary.LittleEndian, &count)
	fmt.Printf("Path count %d bytes: %d\n", unsafe.Sizeof(count), count)
	for i = 0; i < count; i++ {
		s, err := readPath(file)
		fmt.Printf("%+v, %v\n\n", s, err)
	}

	binary.Read(file, binary.LittleEndian, &count)
	fmt.Printf("Shapes count %d bytes: %d\n", unsafe.Sizeof(count), count)
	for i = 0; i < count; i++ {
		s, err := shape.Read(file)
		fmt.Printf("%+v, %v\n\n", s, err)
	}
}
