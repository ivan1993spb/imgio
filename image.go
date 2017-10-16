package imgio

import (
	"errors"
	"image/draw"
	"io"
	"sync"
)

type Image struct {
	img draw.Image
	gen PointsSequenceGenerator
	prw PointReadWriter
	mux sync.RWMutex

	pointCursor uint64
	byteCursor  int
}

func NewImage(img draw.Image, gen PointsSequenceGenerator, prw PointReadWriter) *Image {
	return &Image{
		img: img,
		gen: gen,
		prw: prw,
	}
}

// Read implements io.Reader interface
func (f *Image) Read(p []byte) (n int, err error) {
	f.mux.RLock()
	defer f.mux.RUnlock()

	f.gen.Seek(f.pointCursor)

	if !f.gen.Valid() {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}

	for {
		if !f.gen.Valid() {
			return n, io.EOF
		}
		if n >= len(p) {
			return
		}

		point := f.gen.Current()
		color := f.img.At(point.X, point.Y)
		buff, nBytesRead := f.prw.Read(f.byteCursor, color, point)

		if len(p)-n > nBytesRead {
			copy(p[n:], buff[:nBytesRead])
			n += nBytesRead
			f.gen.Next()
			f.pointCursor++
			f.byteCursor = 0
		} else {
			end := nBytesRead - len(p) + n
			copy(p[n:], buff[:end])
			n += end
			f.byteCursor = end + 1
			return
		}
	}

	return
}

// Write implements io.Writer interface
func (f *Image) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	f.gen.Seek(f.pointCursor)

	for {
		if !f.gen.Valid() {
			return n, errors.New("overflow")
		}
		if len(p) == 0 {
			return n, nil
		}

		point := f.gen.Current()
		srcColor := f.img.At(point.X, point.Y)
		color, writtenBytes := f.prw.Write(p, f.byteCursor, srcColor, point)
		f.img.Set(point.X, point.Y, color)
		f.gen.Next()
		f.pointCursor++
		n += writtenBytes
		p = p[writtenBytes:]
		if f.prw.Size(point) > int64(writtenBytes) {
			f.byteCursor = writtenBytes + 1
		}
	}

	return
}

func (i *Image) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (s *Image) Size() (size int64) {
	s.gen.Rewind()
	for s.gen.Valid() {
		point := s.gen.Current()
		size += s.prw.Size(point)
		s.gen.Next()
	}
	return
}
