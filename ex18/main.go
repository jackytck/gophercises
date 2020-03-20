package main

import (
	"fmt"
	"os/exec"
	"strconv"
)

func main() {
	out, err := primitive("008.jpg", "out.jpg", 80, ellipse)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}

type PrimitiveMode int

const (
	combo PrimitiveMode = iota
	triangle
	rect
	ellipse
	circle
	rotatedrect
	beizers
	rotatedellipse
	polygon
)

func primitive(inputFile, outputFile string, numShapes int, mode PrimitiveMode) (string, error) {
	cmd := exec.Command("primitive", "-i", inputFile, "-o", outputFile, "-n", strconv.Itoa(numShapes), "-m", strconv.Itoa(int(mode)))
	b, err := cmd.CombinedOutput()
	return string(b), err
}
