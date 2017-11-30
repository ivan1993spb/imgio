package imgio

import (
	"image"
	"image/color"
)

type PointReadWriterYCbCr interface {
	Read(start int, c color.YCbCr, p image.Point) ([]byte, int)
	Write(b []byte, start int, src color.YCbCr, p image.Point) (color.YCbCr, int)
	Size(p image.Point) int64
}

type PointReadWriterYCbCrSimple struct {
	Y uint8
}

const PointReadWriterYCbCrSimpleCapacity = 1

func (PointReadWriterYCbCrSimple) Read(start int, c color.YCbCr, p image.Point) ([]byte, int) {
	if start >= PointReadWriterYCbCrSimpleCapacity {
		return []byte{}, 0
	}

	dst := make([]byte, PointReadWriterYCbCrSimpleCapacity-start)
	src := []byte{c.Cb, c.Cr}

	return dst, copy(dst, src[start:])
}

func (prw PointReadWriterYCbCrSimple) Write(b []byte, start int, src color.YCbCr, p image.Point) (color.YCbCr, int) {
	dst := color.YCbCr{src.Y, src.Cb, src.Cr}

	if start >= SimplePoint32Capacity {
		return dst, 0
	}

	dst.Y = prw.Y

	addrs := []*uint8{&dst.Cb, &dst.Cr}

	n := 0
	for i, addr := range addrs[start:] {
		if i >= len(b) {
			break
		}
		*addr = b[i]
		n++
	}

	return dst, n
}

func (PointReadWriterYCbCrSimple) Size(_ image.Point) int64 {
	return PointReadWriterYCbCrSimpleCapacity
}
