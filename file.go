package imgio

import (
	"errors"
	"image/draw"
	"io"
	"sync"
)

type File struct {
	img draw.Image
	gen PointsSequenceGenerator
	prw PointReadWriter
	mux sync.RWMutex

	pointCursor uint64
	byteCursor  int
}

func NewFile(img draw.Image, gen PointsSequenceGenerator, prw PointReadWriter) *File {
	return &File{
		img: img,
		gen: gen,
		prw: prw,
	}
}

// Read implements io.Reader interface
func (f *File) Read(p []byte) (n int, err error) {
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
			break
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
			break
		}
	}

	return
}

// Write implements io.Writer interface
func (f *File) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	f.mux.Lock()
	defer f.mux.Unlock()

	f.gen.Seek(f.pointCursor)

	if !f.gen.Valid() {
		return 0, errors.New("overflow")
	}

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
	}

	return
}

func (s *File) Size() (size int64) {
	s.gen.Rewind()
	for s.gen.Valid() {
		point := s.gen.Current()
		size += s.prw.Size(point)
		s.gen.Next()
	}
	return
}
