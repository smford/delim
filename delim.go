package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	if !term.IsTerminal(0) {
		fmt.Println("Error: not a terminal")
		os.Exit(1)
	}

	width, _, err := term.GetSize(0)

	if err != nil {
		fmt.Println("Error: calculating terminal width")
		os.Exit(2)
	}

	for n := 0; n < width; n++ {
		fmt.Printf("=")
	}
}
