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
	// fmt.Println("allo")

	filename := "res/abydos.hvif"

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
		var styleType style.Type
		binary.Read(file, binary.LittleEndian, &styleType)
		fmt.Printf("Color %d ", i)
		switch styleType {
		case style.StyleSolidColor:
			var c style.SolidColor
			binary.Read(file, binary.LittleEndian, &c)
			fmt.Println(c, "solid")
		case style.StyleGradient:
			var c style.Gradient
			binary.Read(file, binary.LittleEndian, &c)
			fmt.Println(c, "gradient")
		case style.StyleSolidColorNoAlpha:
			var c style.SolidColorNoAlpha
			binary.Read(file, binary.LittleEndian, &c)
			fmt.Println(c, "solid_na")
		case style.StyleSolidGray:
			var c style.SolidGray
			binary.Read(file, binary.LittleEndian, &c)
			fmt.Println(c, "solid_g")
		case style.StyleSolidGrayNoAlpha:
			var c style.SolidGrayNoAlpha
			binary.Read(file, binary.LittleEndian, &c)
			fmt.Println(c, "solid_g_na")
		}
	}
}
