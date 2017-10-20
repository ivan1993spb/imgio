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

const SimplePoint32Capacity = 4

type SimplePoint32ReadWriter struct{}

func (SimplePoint32ReadWriter) Read(start int, c color.Color, p image.Point) ([]byte, int) {
	if start >= SimplePoint32Capacity {
		return []byte{}, 0
	}

	buff := make([]byte, SimplePoint32Capacity-start)
	r, g, b, a := c.RGBA()
	data := []uint32{r, g, b, a}

	n := 0
	for i, v := range data[start:] {
		buff[i] = byte(v)
		n++
	}

	return buff, n
}

func (SimplePoint32ReadWriter) Write(b []byte, start int, src color.Color, p image.Point) (color.Color, int) {
	srcR, srcG, srcB, srcA := src.RGBA()
	c := &color.RGBA{byte(srcR), byte(srcG), byte(srcB), byte(srcA)}

	if start >= SimplePoint32Capacity {
		return c, 0
	}

	addrs := []*uint8{&c.R, &c.G, &c.B, &c.A}

	n := 0
	for i, addr := range addrs[start:] {
		if i >= len(b) {
			break
		}
		*addr = b[i]
		n++
	}

	return c, n
}

func (SimplePoint32ReadWriter) Size(_ image.Point) int64 {
	return SimplePoint32Capacity
}

const SimplePoint64Capacity = 8

type SimplePoint64ReadWriter struct{}

func (SimplePoint64ReadWriter) Read(start int, c color.Color, p image.Point) ([]byte, int) {
	if start >= SimplePoint64Capacity {
		return []byte{}, 0
	}

	buff := make([]byte, SimplePoint64Capacity-start)
	r, g, b, a := c.RGBA()
	data := []uint32{r, g, b, a}
	n := 0
	skip := start % 2
	for _, v := range data[start/2:] {
		if skip == 1 {
			skip = 0
		} else {
			buff[n] = byte(v >> 8)
			n++
		}
		buff[n] = byte(v)
		n++
	}

	return buff, n
}

func (SimplePoint64ReadWriter) Write(b []byte, start int, src color.Color, p image.Point) (color.Color, int) {
	srcR, srcG, srcB, srcA := src.RGBA()
	c := &color.RGBA64{uint16(srcR), uint16(srcG), uint16(srcB), uint16(srcA)}

	if start >= SimplePoint64Capacity {
		return c, 0
	}

	addrs := []*uint16{&c.R, &c.G, &c.B, &c.A}
	n := 0
	i := 0
	skip := start % 2
	for _, addr := range addrs[start/2:] {
		if i >= len(b) {
			break
		}
		if skip == 1 {
			skip = 0
			i++

			*addr &= 0xff00
			*addr |= uint16(b[i])
			n++
			i++
		} else {
			*addr &= 0x00ff
			*addr |= uint16(b[i]) << 8
			n++
			i++

			if i >= len(b) {
				return c, n
			}

			*addr &= 0xff00
			*addr |= uint16(b[i])
			n++
			i++
		}
	}

	return c, n
}

func (SimplePoint64ReadWriter) Size(_ image.Point) int64 {
	return SimplePoint64Capacity
}

type SmartPoint8ReadWriter struct{}

const SmartPoint8Capacity = 1

func (SmartPoint8ReadWriter) Read(start int, c color.Color, p image.Point) ([]byte, int) {
	// TODO: Implement method.
	return []byte{}, 0
}

func (SmartPoint8ReadWriter) Write(b []byte, start int, src color.Color, p image.Point) (color.Color, int) {
	// TODO: Implement method.
	return nil, 0
}

func (SmartPoint8ReadWriter) Size(_ image.Point) int64 {
	return SmartPoint8Capacity
}

type GentlePoint16ReadWriter struct{}

const GentlePoint16Capacity = 2

func (GentlePoint16ReadWriter) Read(start int, c color.Color, p image.Point) ([]byte, int) {
	if start >= GentlePoint16Capacity {
		return []byte{}, 0
	}

	buff := make([]byte, GentlePoint16Capacity-start)
	r, g, b, a := c.RGBA()
	n := 0

	if start == 0 {
		buff[n] = byte(r) & 0x0f << 4
		buff[n] |= byte(g) & 0x0f
		n++
	}

	buff[n] = byte(b) & 0x0f << 4
	buff[n] |= byte(a) & 0x0f
	n++

	return buff, n
}

func (GentlePoint16ReadWriter) Write(b []byte, start int, src color.Color, p image.Point) (color.Color, int) {
	srcR, srcG, srcB, srcA := src.RGBA()
	c := &color.RGBA{uint8(srcR), uint8(srcG), uint8(srcB), uint8(srcA)}

	if len(b) == 0 {
		return c, 0
	}

	if start >= GentlePoint16Capacity {
		return c, 0
	}

	i := 0

	if start == 0 {
		c.R &= 0xf0
		c.R |= b[i] & 0xf0 >> 4
		c.G &= 0xf0
		c.G |= b[i] & 0x0f
		i++
		c.B &= 0xf0
		c.B |= b[i] & 0xf0 >> 4
		c.A &= 0xf0
		c.A |= b[i] & 0x0f
		i++
	} else if start == 1 {
		c.B &= 0xf0
		c.B |= b[i] & 0xf0 >> 4
		c.A &= 0xf0
		c.A |= b[i] & 0x0f
		i++
	}

	return c, i
}

func (GentlePoint16ReadWriter) Size(_ image.Point) int64 {
	return GentlePoint16Capacity
}
