package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gucio321/cimguizer"
)

func main() {
	path := flag.String("path", "", "Path the .h file")
	flag.Parse()

	if len(*path) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	fileData, err := os.ReadFile(*path)
	if err != nil {
		log.Fatal(err)
	}

	c, err := cimguizer.Parse(fileData)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(c)
	fmt.Println(c.StructAndEnums())
}
