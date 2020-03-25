package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

// Mode defines the shapes used when transforming images.
type Mode int

// Modes supported by the primitive package.
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedrect
	ModeBeizers
	ModeRotatedellipse
	ModePolygon
)

// WithMode is an option for the Transform function that will define the
// mode you want to use. By default, ModeTriangle will be used.
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply a primitive
// transformation to it, then return a reader to the resulting image.
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}
	in, err := tempfile("in_", ext)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	defer os.Remove(in.Name())

	out, err := tempfile("out_", ext)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	defer os.Remove(out.Name())

	// read image into in file
	_, err = io.Copy(in, image)
	if err != nil {
		return nil, errors.New("primitive: failed to copy image into temp input file")
	}

	// run primitive w/ -i in.Name() -o out.Name()
	stdCombo, err := run(in.Name(), out.Name(), numShapes, args...)
	if err != nil {
		return nil, fmt.Errorf("primitive: failed to run the primitive command. stdcombo=%s", stdCombo)
	}

	// read out into a reader, return reader, delete out
	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, out)
	if err != nil {
		return nil, errors.New("primitive: failed to copy output file into byte buffer")
	}

	return b, nil
}

func run(inputFile, outputFile string, numShapes int, args ...string) (string, error) {
	a := []string{"-i", inputFile, "-o", outputFile, "-n", strconv.Itoa(numShapes)}
	args = append(a, args...)
	cmd := exec.Command("primitive", args...)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func tempfile(prefix, ext string) (*os.File, error) {
	f, err := ioutil.TempFile("", prefix)
	if err != nil {
		return nil, errors.New("primitive: failed to create temporary file")
	}
	defer f.Close()
	defer os.Remove(f.Name())
	return os.Create(fmt.Sprintf("%s.%s", f.Name(), ext))

}
