package graphicsmagick

import "testing"

func TestReadImage(t *testing.T) {
	im, err := ReadImage("fixtures/image.jpg")
	if im == nil {
		t.Fatalf("couldn't open existing file - err: %v", err)
	}
	defer im.Destroy()

	im, err = ReadImage("fixtures/!!missing!!")
	if im != nil {
		t.Fatalf("im present for missing file")
	}

	im, err = ReadImage("fixtures/doc.pdf[0]")
	if im == nil {
		t.Fatalf("couldn't open pdf[0] - err: %v", err)
	}
	defer im.Destroy()
}
