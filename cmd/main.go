package main

import (
	"fmt"
	"log"
	"os"

	"hvif"
)

func main() {
	filename := "testdata/ime.hvif"

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer file.Close()

	img, err := hvif.ReadImage(file)
	fmt.Println(img, err)
}
