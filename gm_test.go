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
	require.Nil(t, err)
	defer im.Destroy()

	// missing
	im, err = ReadImage("fixtures/!!missing!!")
	require.Nil(t, im)

	// PDF
	im, err = ReadImage("fixtures/doc.pdf[0]")
	require.Nil(t, err)
	defer im.Destroy()

	// resize
	im, err = ReadImage("fixtures/image.jpg")
	require.Nil(t, err)
	defer im.Destroy()
	res, err := im.Resize(10, 10, "", 1)
	require.Nil(t, err)
	defer res.Destroy()
	require.Equal(t, uint(10), res.Columns())
	require.Equal(t, uint(10), res.Rows())

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
	require.Equal(t, q, iminfo.Quality())

	// background color
	color, err := QueryColorDatabase("red")
	require.Nil(t, err)
	iminfo.SetBackgroundColor(color)
	require.Equal(t, color.Hex(), iminfo.BackgroundColor().Hex())

	// filename
	filename := "xxx"
	iminfo.SetFilename(filename)
	require.Equal(t, filename, iminfo.Filename())

	// WriteImage
	f, err := ioutil.TempFile("", "WriteImage-test")
	require.Nil(t, err)
	defer func() {
		os.Remove(f.Name())
	}()
	f.Close()
	im, err := ReadImage("fixtures/image.jpg")
	require.Nil(t, err)
	defer im.Destroy()
	iminfo = NewImageInfo()
	defer iminfo.Destroy()
	im.SetFilename(f.Name())
	err = iminfo.WriteImage(im)
	require.Nil(t, err)
}

func TestPixelPacket(t *testing.T) {
	// QueryColorDatabase
	name := "blue"
	color, err := QueryColorDatabase(name)
	require.Nil(t, err)
	require.Equal(t, "0000ff00", color.Hex())

	// QueryColorDatabase: unknown
	name = "!!unknown!!"
	color, err = QueryColorDatabase(name)
	require.NotNil(t, err)
}

func TestImage(t *testing.T) {
	// filename
	im := AllocateImage()
	require.Equal(t, "", im.Filename())
	im.SetFilename("abc")
	require.Equal(t, "abc", im.Filename())
}
