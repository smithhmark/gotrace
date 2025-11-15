package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/smithhmark/gotracer/internal/renderer"
)

func main() {

	scene := renderer.Render()

	filename := "test.png"
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("failed to open %q, quiting\n", filename)
		os.Exit(1)
	}
	defer outputFile.Close()

	if err := png.Encode(outputFile, scene); err != nil {
		fmt.Println("failed to save image to file")
		os.Exit(1)
	}
	fmt.Println("Saved image to file")
}
