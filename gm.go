package graphicsmagick

/*
#cgo pkg-config: GraphicsMagick
#cgo LDFLAGS: -DQuantumDepth=8
#include <magick/api.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	gmTrue  = 1
	gmFalse = 0
)

func init() {
	C.InitializeMagick(nil)
	// TODO DestroyMagick() ?
}

func ReadImage(path string) (*Image, error) {
	iminfo := NewImageInfo()
	defer iminfo.Destroy()
	iminfo.SetFilename(path)
	return iminfo.ReadImage()
}

type ImageInfo struct {
	c *C.ImageInfo
}

func NewImageInfo() *ImageInfo {
	cinf := C.CloneImageInfo(nil)
	return &ImageInfo{cinf}
}

func (inf *ImageInfo) Destroy() {
	defer C.DestroyImageInfo(inf.c)
}

func (inf *ImageInfo) SetFilename(filename string) {
	gmStrcpy(&inf.c.filename, filename)
}

func (inf *ImageInfo) Filename() string {
	return gmGoString(inf.c.filename)
}

func (inf *ImageInfo) SetQuality(quality uint) {
	inf.c.quality = C.ulong(quality)
}

func (inf *ImageInfo) Quality() uint {
	return uint(inf.c.quality)
}

func (inf *ImageInfo) ReadImage() (*Image, error) {
	exc := newExceptionInfo()
	defer exc.Destroy()
	cim := C.ReadImage(inf.c, exc.c)
	if cim == nil {
		return nil, exc.MustError("while reading file")
	}
	return &Image{cim}, nil
}

type exceptionInfo struct {
	c *C.ExceptionInfo
}

func newExceptionInfo() *exceptionInfo {
	cexc := &C.ExceptionInfo{}
	C.GetExceptionInfo(cexc)
	return &exceptionInfo{cexc}
}

func (exc *exceptionInfo) Destroy() {
	C.DestroyExceptionInfo(exc.c)
}

// Flow from GM CatchException()
func (exc *exceptionInfo) GetError() error {
	if exc.c.signature != C.MagickSignature {
		return nil
	}
	if exc.c.severity == C.UndefinedException {
		return nil
	}
	// TODO extract severity? (ExceptionType)
	reason := C.GoString(exc.c.reason)
	description := C.GoString(exc.c.description)
	return errors.New(reason + ": " + description)
}

func (exc *exceptionInfo) MustError(msg string) error {
	err := exc.GetError()
	if err != nil {
		msg += ": " + err.Error()
	} else {
		msg += ": unknown error"
	}
	return errors.New(msg)
}

type Image struct {
	c *C.Image
}

func (im *Image) Destroy() {
	C.DestroyImage(im.c)
}

type PixelPacket struct {
	c C.PixelPacket
}

func QueryColorDatabase(name string) (*PixelPacket, error) {
	cpxpacket := C.PixelPacket{}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	exc := newExceptionInfo()
	defer exc.Destroy()
	res := C.QueryColorDatabase(cname, &cpxpacket, exc.c)
	if res != gmTrue {
		return nil, exc.MustError("in QueryColorDatabase()")
	}
	return &PixelPacket{cpxpacket}, nil
}

func (pp *PixelPacket) Red() uint8 {
	return uint8(pp.c.red)
}

func (pp *PixelPacket) Green() uint8 {
	return uint8(pp.c.green)
}

func (pp *PixelPacket) Blue() uint8 {
	return uint8(pp.c.blue)
}

func (pp *PixelPacket) Opacity() uint8 {
	return uint8(pp.c.opacity)
}

func (pp *PixelPacket) Hex() string {
	return fmt.Sprintf("%02x%02x%02x%02x",
		pp.Red(), pp.Green(), pp.Blue(), pp.Opacity())
}

func gmGoString(str [C.MaxTextExtent]C.char) string {
	gostr := make([]byte, 0, len(str))
	nullChar := C.char(0)
	for _, c := range str {
		if c == nullChar {
			break
		}
		gostr = append(gostr, byte(c))
	}
	return string(gostr)
}

func gmStrcpy(dst *[C.MaxTextExtent]C.char, src string) {
	nsrc, ndst := len(src), len(dst)
	max := ndst - 1
	if max > nsrc {
		max = nsrc
	}
	src = src[0:max]
	for i, c := range src {
		dst[i] = C.char(c)
	}
	dst[len(src)] = 0
}
