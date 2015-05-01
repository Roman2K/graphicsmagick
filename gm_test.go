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
		t.Fatalf("SetQuality set %v instead of %v", actual, q)
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

	name = "!!missing!!"
	color, err = QueryColorDatabase(name)
	if err == nil {
		t.Fatalf("QueryColorDatabase found missing color %s", name)
	}
}
