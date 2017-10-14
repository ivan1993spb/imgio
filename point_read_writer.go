package imgio

import (
	"errors"
	"image/color"
)

type PointReadWriter interface {
	Read(c color.Color) []byte
	Write(b []byte) (color.Color, error)
}

type SimplePointReadWriter struct{}

func (SimplePointReadWriter) Read(c color.Color) []byte {
	buff := make([]byte, 4)
	r, g, b, a := c.RGBA()

	buff[0] = byte(r)
	buff[1] = byte(g)
	buff[2] = byte(b)
	buff[3] = byte(a)

	return buff
}

func (SimplePointReadWriter) Write(b []byte) (color.Color, error) {
	c := &color.RGBA{}

	if len(b) > 0 {
		c.R = b[0]
	} else {
		return c, nil
	}

	if len(b) > 1 {
		c.G = b[1]
	} else {
		return c, nil
	}

	if len(b) > 2 {
		c.B = b[2]
	} else {
		return c, nil
	}

	if len(b) > 3 {
		c.A = b[3]
	} else {
		return c, nil
	}

	if len(b) > 4 {
		return c, errors.New("overflow")
	}

	return c, nil
}
