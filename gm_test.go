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

func TestImageInfo(t *testing.T) {
	iminfo := NewImageInfo()
	defer iminfo.Destroy()
	q := uint(34)
	iminfo.SetQuality(q)
	if actual := iminfo.Quality(); actual != 34 {
		t.Fatalf("SetQuality failed: actual = %v, expected = %v", actual, q)
	}
}
