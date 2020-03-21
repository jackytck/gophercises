package main

import (
	"io"
	"os"

	"github.com/jackytck/gophercises/ex18/primitive"
)

func main() {
	inFile, err := os.Open("008.jpg")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	out, err := primitive.Transform(inFile, 50)
	if err != nil {
		panic(err)
	}
	err = os.Remove("out.jpg")
	if err != nil {
		panic(err)
	}
	outFile, err := os.Create("out.jpg")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, out)
	if err != nil {
		panic(err)
	}
}
