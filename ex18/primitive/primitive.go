package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Mode defines the shapes used when transforming images.
type Mode int

// Modes supported by the primitive package
const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

// WithMode is an option for the Transform function tha will define the
// mode you want to use. By defaukt, ModeTriangle will be used.
func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-n", fmt.Sprintf("%d", mode)}
	}
}

// Transform will take the provided image and apply a primitive
// transformation to it, then return a reader to the resulting image.
func Transform(image io.Reader, ext string, numShapes int, opts ...func() []string) (io.Reader, error) {
	in, err := tempfile("in_", ext)
	if err != nil {
		return nil, err
	}

	defer os.Remove(in.Name())

	out, err := tempfile("out_", ext)
	if err != nil {
		return nil, err
	}

	defer os.Remove(out.Name())

	_, err = io.Copy(in, image)

	if err != nil {
		return nil, err
	}

	stdCombo, err := primitive(in.Name(), out.Name(), numShapes, ModeCombo)
	if err != nil {
		return nil, err
	}
	fmt.Println(stdCombo)

	b := bytes.NewBuffer(nil)

	_, err = io.Copy(b, out)

	if err != nil {
		return nil, err
	}

	return b, nil

}

func primitive(inputFile, outputFile string, numShapes int, mode Mode) (string, error) {
	args := fmt.Sprintf("-i %s -o %s -n %d -m %d", inputFile, outputFile, numShapes, mode)
	cmd := exec.Command("primitive", strings.Fields(args)...)
	b, err := cmd.CombinedOutput()

	return string(b), err
}

func tempfile(prefix, ext string) (*os.File, error) {
	in, err := ioutil.TempFile("", "in_")
	if err != nil {
		return nil, errors.New("tempfile: failed to create temporary file")
	}

	defer os.Remove(in.Name())
	return os.Create(fmt.Sprintf("%s.%s", in.Name(), ext))

}
