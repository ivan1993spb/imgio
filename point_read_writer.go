package imgio

import (
	"image"
	"image/color"
)

type PointReadWriter interface {
	// Read reads bytes from color c from position start on point p
	Read(start int, c color.Color, p image.Point) ([]byte, int)
	// Write writes bytes b into color c starts on position start on point p and returns number of written bytes
	Write(b []byte, start int, c color.Color, p image.Point) (color.Color, int)
	// Return number of bytes possible to be written to point p on current image
	Size(p image.Point) int64
}

const SimplePointCapacity = 4

type SimplePointReadWriter struct{}

func (SimplePointReadWriter) Read(start int, c color.Color, p image.Point) ([]byte, int) {
	if start >= SimplePointCapacity {
		return []byte{}, 0
	}

	buff := make([]byte, SimplePointCapacity-start)
	r, g, b, a := c.RGBA()
	data := []uint32{r, g, b, a}

	n := 0
	for i, v := range data[start:] {
		buff[i] = byte(v)
		n++
	}

	return buff, n
}

func (SimplePointReadWriter) Write(b []byte, start int, src color.Color, p image.Point) (color.Color, int) {
	srcR, srcG, srcB, srcA := src.RGBA()
	c := &color.RGBA{byte(srcR), byte(srcG), byte(srcB), byte(srcA)}

	if start >= SimplePointCapacity {
		return c, 0
	}

	addrs := []*uint8{&c.R, &c.G, &c.B, &c.A}

	n := 0
	for i, addr := range addrs[start:] {
		*addr = b[i]
		n++
	}

	return c, n
}

func (SimplePointReadWriter) Size(_ image.Point) int64 {
	return SimplePointCapacity
}
