package graphicsmagick

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TODO use testify where relevant

func TestReadImage(t *testing.T) {
	// normal
	im, err := ReadImage("fixtures/image.jpg")
	if im == nil {
		t.Fatalf("couldn't open existing file - err: %v", err)
	}
	defer im.Destroy()

	// missing
	im, err = ReadImage("fixtures/!!missing!!")
	if im != nil {
		t.Fatalf("im present for missing file")
	}

	// PDF
	im, err = ReadImage("fixtures/doc.pdf[0]")
	if im == nil {
		t.Fatalf("couldn't open pdf[0] - err: %v", err)
	}
	defer im.Destroy()

	// resize
	im, err = ReadImage("fixtures/image.jpg")
	if err != nil {
		t.Fatalf("in ReadImage(): %v", err)
	}
	defer im.Destroy()
	res, err := im.Resize(10, 10, "", 1)
	if err != nil {
		t.Fatalf("Resize() failed: %v", err)
	}
	defer res.Destroy()
	if x, y := res.Columns(), res.Rows(); x != 10 || y != 10 {
		t.Fatalf("incorrect resize: %dx%d instead of 10x10", x, y)
	}

	// resize: unknown filter
	_, err = im.Resize(10, 10, "!!unknown!!", 1)
	if err == nil || !strings.Contains(err.Error(), "unknown filter") {
		t.Fatalf("didn't return an error for unknown filter")
	}
}

func TestImageInfo(t *testing.T) {
	iminfo := NewImageInfo()
	defer iminfo.Destroy()

	// quality
	q := uint(34)
	iminfo.SetQuality(q)
	if actual := iminfo.Quality(); actual != 34 {
		t.Fatalf("SetQuality set %v instead of %v", actual, q)
	}

	// background color
	color, err := QueryColorDatabase("red")
	if err != nil {
		t.Fatalf("QueryColorDatabase() failed: %v", err)
	}
	iminfo.SetBackgroundColor(color)
	if actual, expected := iminfo.BackgroundColor().Hex(), color.Hex(); actual != expected {
		t.Fatalf("SetBackground set %v instead of %v", actual, expected)
	}

	// filename
	filename := "xxx"
	iminfo.SetFilename(filename)
	if actual := iminfo.Filename(); actual != filename {
		t.Fatalf("SetFilename set %v instead of %v", actual, filename)
	}

	// WriteImage
	f, err := ioutil.TempFile("", "WriteImage-test")
	if err != nil {
		t.Fatalf("TempFile error: %v", err)
	}
	defer func() {
		os.Remove(f.Name())
	}()
	f.Close()
	im, err := ReadImage("fixtures/image.jpg")
	if err != nil {
		t.Fatalf("ReadImage() error: %v", err)
	}
	defer im.Destroy()
	iminfo = NewImageInfo()
	defer iminfo.Destroy()
	im.SetFilename(f.Name())
	err = iminfo.WriteImage(im)
	if err != nil {
		t.Fatalf("WriteImage() error: %v", err)
	}
}

func TestPixelPacket(t *testing.T) {
	name := "blue"
	color, err := QueryColorDatabase(name)
	if err != nil {
		t.Fatalf("QueryColorDatabase couldn't find %s: %v", name, err)
	}
	if actual, expected := color.Hex(), "0000ff00"; actual != expected {
		t.Fatalf("returned color is %s instead of %s", actual, expected)
	}

	name = "!!unknown!!"
	color, err = QueryColorDatabase(name)
	if err == nil {
		t.Fatalf("QueryColorDatabase found missing color %s", name)
	}
}

func TestImage(t *testing.T) {
	// filename
	im := AllocateImage()
	require.Equal(t, "", im.Filename())
	im.SetFilename("abc")
	require.Equal(t, "abc", im.Filename())
}
