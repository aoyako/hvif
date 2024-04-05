package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	filename := "res/ime"

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer file.Close()

	img, err := ReadImage(file)
	fmt.Println(img, err)
}
