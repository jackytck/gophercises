package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	filename := flag.String("html", "", "a html file")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}

	links, err := Parse(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(links)
}
