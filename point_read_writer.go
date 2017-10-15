package imgio

import (
	"errors"
	"image"
	"image/color"
)

type PointReadWriter interface {
	// Read reads bytes from color c on point p
	Read(c color.Color, p image.Point) ([]byte, int)
	// Write writes bytes b into color c on point p and returns number of written bytes
	Write(b []byte, c color.Color, p image.Point) (color.Color, int, error)
	// Return number of bytes possible to be written to point p on current image
	Size(p image.Point) int64
}

type SimplePointReadWriter struct{}

func (SimplePointReadWriter) Read(c color.Color, p image.Point) ([]byte, int) {
	buff := make([]byte, 4)
	r, g, b, a := c.RGBA()

	buff[0] = byte(r)
	buff[1] = byte(g)
	buff[2] = byte(b)
	buff[3] = byte(a)

	return buff, 4
}

func (SimplePointReadWriter) Write(b []byte, _ color.Color, _ image.Point) (color.Color, int, error) {
	c := &color.RGBA{}

	if len(b) > 0 {
		c.R = b[0]
	} else {
		return c, 0, nil
	}

	if len(b) > 1 {
		c.G = b[1]
	} else {
		return c, 1, nil
	}

	if len(b) > 2 {
		c.B = b[2]
	} else {
		return c, 2, nil
	}

	if len(b) > 3 {
		c.A = b[3]
	} else {
		return c, 3, nil
	}

	if len(b) > 4 {
		return c, 4, errors.New("overflow")
	}

	return c, 4, nil
}

func (SimplePointReadWriter) Size(p image.Point) int64 {
	return 4
}
