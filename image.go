package imgio

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io"
	"sync"
)

type Image struct {
	img draw.Image
	gen PointsSequenceGenerator
	prw PointReadWriter
	mux sync.RWMutex

	byteCursor int
}

func NewImage(img draw.Image, gen PointsSequenceGenerator, prw PointReadWriter) *Image {
	return &Image{
		img: img,
		gen: gen,
		prw: prw,
	}
}

// Read implements io.Reader interface
func (i *Image) Read(p []byte) (n int, err error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if !i.gen.Valid() {
		return 0, io.EOF
	}

	if len(p) == 0 {
		return 0, nil
	}

	for {
		if !i.gen.Valid() {
			return n, io.EOF
		}
		if n >= len(p) {
			return
		}

		point := i.gen.Current()
		color := i.img.At(point.X, point.Y)
		buff, nBytesRead := i.prw.Read(i.byteCursor, color, point)

		if len(p)-n >= nBytesRead {
			copy(p[n:], buff[:nBytesRead])
			n += nBytesRead
			i.gen.Next()
			i.byteCursor = 0
		} else {
			end := nBytesRead - len(p) + n
			copy(p[n:], buff[:end])
			n += end
			i.byteCursor = end + 1
			return
		}
	}

	return
}

var ErrOverflow = errors.New("Overflow")

// Write implements io.Writer interface
func (i *Image) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	i.mux.Lock()
	defer i.mux.Unlock()

	if !i.gen.Valid() {
		return 0, ErrOverflow
	}

	for {
		if len(p) == 0 {
			return n, nil
		}
		if !i.gen.Valid() {
			return n, ErrOverflow
		}

		point := i.gen.Current()
		srcColor := i.img.At(point.X, point.Y)
		color, writtenBytes := i.prw.Write(p, i.byteCursor, srcColor, point)
		i.img.Set(point.X, point.Y, color)
		i.gen.Next()
		n += writtenBytes
		p = p[writtenBytes:]
		if i.prw.Size(point) > int64(writtenBytes) {
			i.byteCursor = writtenBytes + 1
		}
	}

	return
}

func (i *Image) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (i *Image) Size() (size int64) {
	i.gen.Rewind()
	defer i.gen.Rewind()
	for i.gen.Valid() {
		point := i.gen.Current()
		size += i.prw.Size(point)
		i.gen.Next()
	}
	return
}

// ColorModel implements image.Image interface
func (i *Image) ColorModel() color.Model {
	return i.img.ColorModel()
}

// Bounds implements image.Image interface
func (i *Image) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// At  implements image.Image interface
func (i *Image) At(x, y int) color.Color {
	return i.img.At(x, y)
}
