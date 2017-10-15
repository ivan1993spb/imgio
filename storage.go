package imgio

import (
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type Storage struct {
	img draw.Image
	gen PointsSequenceGenerator
	prw PointReadWriter
}

func NewStorage(img draw.Image, gen PointsSequenceGenerator, prw PointReadWriter) *Storage {
	return &Storage{
		img: img,
		gen: gen,
		prw: prw,
	}
}

// Read implements io.Reader interface
func (s *Storage) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	cursor := 0

	s.gen.Rewind()
	for s.gen.Valid() && cursor < len(p) {
		point := s.gen.Current()
		c := s.img.At(point.X, point.Y)
		b, n := s.prw.Read(c, point)
		copy(p[cursor:cursor+n], b[:n])
		cursor += n
		s.gen.Next()
	}

	return cursor, io.EOF
}

// Write implements io.Writer interface
func (s *Storage) Write(p []byte) (n int, err error) {
	s.gen.Rewind()
	for s.gen.Valid() && len(p) > 0 {
		point := s.gen.Current()
		c := s.img.At(point.X, point.Y)
		color, written, _ := s.prw.Write(p, c, point)
		s.img.Set(point.X, point.Y, color)
		s.gen.Next()
		n += written
		p = p[written:]
	}

	return
}

func (s *Storage) Size() (size int64) {
	s.gen.Rewind()
	for s.gen.Valid() {
		point := s.gen.Current()
		size += s.prw.Size(point)
		s.gen.Next()
	}
	return
}
