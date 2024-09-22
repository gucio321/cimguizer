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
	funcsOutputPath := flag.String("funcs", "", "Path to the output file with functions (definitions.json)")
	samOutputPath := flag.String("sam", "", "Path to Structs And Enums json output (structs_and_enums.json)")
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

	// write funcs
	if len(*funcsOutputPath) > 0 {
		funcs, err := c.Funcs()
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(*funcsOutputPath, []byte(funcs), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Functions saved to", *funcsOutputPath)
	}

	// write structs and enums
	if len(*samOutputPath) > 0 {
		sam, err := c.StructAndEnums()
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(*samOutputPath, []byte(sam), 0644)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Structs and Enums saved to", *samOutputPath)
	}

	sam, _ := (c.StructAndEnums())
	fmt.Println("----")
	fmt.Println(sam)
}
