package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"
)

func main() {
	filename := "res/ime"

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer file.Close()

	magic := make([]byte, 4)
	nread, err := io.ReadFull(file, magic)
	if err != nil {
		log.Println("Error reading file:", err)
	}

	fmt.Printf("Magic %d bytes: %s\n", nread, magic)

	var count byte
	binary.Read(file, binary.LittleEndian, &count)
	fmt.Printf("Color count %d bytes: %d\n", unsafe.Sizeof(count), count)

	var i byte
	for i = 0; i < 1; i++ {
		s, err := readStyle(file)
		fmt.Printf("%+v, %v\n\n", s, err)
		fmt.Println(&Color{Red: 1, Green: 1, Blue: 1, Alpha: 116} == s)
		// fmt.Printf("%T, %p, %v \n\n", s, &s, err)
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
		s, err := readShape(file)
		fmt.Printf("%+v, %v\n\n", s, err)
	}
}
