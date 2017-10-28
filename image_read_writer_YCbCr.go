package imgio

import (
	"errors"
	"image"
	"image/color"
	"io"
	"sync"
)

type ImageReadWriterYCbCr struct {
	img *image.YCbCr
	gen PointsSequenceGenerator
	prw PointReadWriterYCbCr
	mux sync.RWMutex

	byteCursor int
}

func NewImageReadWriterYCbCr(img *image.YCbCr, gen PointsSequenceGenerator, prw PointReadWriterYCbCr) *ImageReadWriterYCbCr {
	return &ImageReadWriterYCbCr{
		img: img,
		gen: gen,
		prw: prw,
	}
}

// Read implements io.Reader interface
func (i *ImageReadWriterYCbCr) Read(p []byte) (n int, err error) {
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
		buff, nBytesRead := i.prw.Read(i.byteCursor, i.img.YCbCrAt(point.X, point.Y), point)

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

var ErrImageReadWriterYCbCrOverflow = errors.New("Overflow")

// Write implements io.Writer interface
func (i *ImageReadWriterYCbCr) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	i.mux.Lock()
	defer i.mux.Unlock()

	if !i.gen.Valid() {
		return 0, ErrImageReadWriterYCbCrOverflow
	}

	for {
		if len(p) == 0 {
			return n, nil
		}
		if !i.gen.Valid() {
			return n, ErrImageReadWriterYCbCrOverflow
		}

		point := i.gen.Current()
		srcColor := i.img.YCbCrAt(point.X, point.Y)
		c, writtenBytes := i.prw.Write(p, i.byteCursor, srcColor, point)
		i.img.Y[i.img.YOffset(point.X, point.Y)] = c.Y
		i.img.Cb[i.img.COffset(point.X, point.Y)] = c.Cb
		i.img.Cr[i.img.COffset(point.X, point.Y)] = c.Cr

		i.gen.Next()
		n += writtenBytes
		p = p[writtenBytes:]
		if i.prw.Size(point) > int64(writtenBytes) {
			i.byteCursor = writtenBytes + 1
		}
	}

	return
}

func (i *ImageReadWriterYCbCr) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (i *ImageReadWriterYCbCr) Size() (size int64) {
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
func (i *ImageReadWriterYCbCr) ColorModel() color.Model {
	return i.img.ColorModel()
}

// Bounds implements image.Image interface
func (i *ImageReadWriterYCbCr) Bounds() image.Rectangle {
	return i.img.Bounds()
}

// At  implements image.Image interface
func (i *ImageReadWriterYCbCr) At(x, y int) color.Color {
	return i.img.At(x, y)
}
