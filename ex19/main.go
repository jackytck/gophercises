package main

import (
	"math/rand"
	"os"
	"time"

	svg "github.com/ajstarks/svgo"
)

func rn(n int) int { return rand.Intn(n) }

func main() {
	f, err := os.OpenFile("demo.svg", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	canvas := svg.New(f)
	width := 500
	height := 500
	nstars := 250
	style := "font-size:48pt;fill:white;text-anchor:middle"

	rand.Seed(time.Now().Unix())
	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height)
	for i := 0; i < nstars; i++ {
		canvas.Circle(rn(width), rn(height), rn(3), "fill:white")
	}
	canvas.Circle(width/2, height, width/2, "fill:rgb(77, 117, 232)")
	canvas.Text(width/2, height*4/5, "hello, world", style)
	canvas.End()
}
