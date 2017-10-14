package imgio

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Storage struct {
	img image.Image
	gen PointsSequenceGenerator
	prw PointReadWriter
}

func NewStorage(img image.Image, gen PointsSequenceGenerator, prw PointReadWriter) *Storage {
	return &Storage{img, gen, prw}
}

func (s *Storage) Read(b []byte) (int64, error) {
	return 0, nil
}
