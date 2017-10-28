package imgio

import (
	"image"
	"image/color"
)

type PointReadWriterRGBA interface {
	Read(start int, c color.RGBA, p image.Point) ([]byte, int)
	Write(b []byte, start int, src color.RGBA, p image.Point) (color.RGBA, int)
	Size(p image.Point) int64
}

type PointReadWriterRGBASimple struct{}

const PointReadWriterRGBASimpleCapacity = 1

func (PointReadWriterRGBASimple) Read(start int, c color.RGBA, p image.Point) ([]byte, int) {
	// TODO: Implement method.
	return []byte{}, 0
}

func (PointReadWriterRGBASimple) Write(b []byte, start int, src color.RGBA, p image.Point) (color.RGBA, int) {
	// TODO: Implement method.
	return color.RGBA{}, 0
}

func (PointReadWriterRGBASimple) Size(_ image.Point) int64 {
	return PointReadWriterRGBASimpleCapacity
}
