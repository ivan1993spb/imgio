package imgio

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimplePoint32ReadWriter_Read(t *testing.T) {
	tests := []struct {
		startOn int
		color   color.RGBA
		point   image.Point

		expectedBuff   []byte
		expectedNumber int
	}{
		{0, color.RGBA{'a', 'b', 'c', 'd'}, image.Point{}, []byte{'a', 'b', 'c', 'd'}, 4},
		{1, color.RGBA{R: 'a'}, image.Point{}, []byte{0, 0, 0}, 3},
		{2, color.RGBA{B: 'a', A: 'b'}, image.Point{}, []byte{'a', 'b'}, 2},
		{4, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, []byte{}, 0},
	}

	for i, test := range tests {
		b, n := SimplePoint32ReadWriter{}.Read(test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedBuff, b, "Test index %d", i)
	}
}

func Test_SimplePoint32ReadWriter_Write(t *testing.T) {
	tests := []struct {
		buff    []byte
		startOn int
		color   color.RGBA
		point   image.Point

		expectedColor  *color.RGBA
		expectedNumber int
	}{
		{[]byte{'a', 'b', 'c', 'd'}, 0, color.RGBA{}, image.Point{}, &color.RGBA{'a', 'b', 'c', 'd'}, 4},
		{[]byte{'a', 'b', 'c', 'd'}, 1, color.RGBA{R: 'a'}, image.Point{}, &color.RGBA{'a', 'a', 'b', 'c'}, 3},
		{[]byte{'a', 'b', 'c', 'd'}, 2, color.RGBA{R: 'a'}, image.Point{}, &color.RGBA{'a', 0, 'a', 'b'}, 2},
		{[]byte{'a', 'b', 'c', 'd'}, 2, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA{'e', 'f', 'a', 'b'}, 2},
		{[]byte{'a', 'b', 'c', 'd'}, 4, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA{'e', 'f', 'g', 'h'}, 0},
		{[]byte{'a', 'b', 'c', 'd', 'e', 'f'}, 0, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA{'a', 'b', 'c', 'd'}, 4},
		{[]byte{'a'}, 0, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA{'a', 'f', 'g', 'h'}, 1},
		{[]byte{'i'}, 2, color.RGBA{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA{'e', 'f', 'i', 'h'}, 1},
	}

	for i, test := range tests {
		c, n := SimplePoint32ReadWriter{}.Write(test.buff, test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedColor, c, "Test index %d", i)
	}
}

func Test_SimplePoint64ReadWriter_Read(t *testing.T) {
	tests := []struct {
		startOn int
		color   color.RGBA64
		point   image.Point

		expectedBuff   []byte
		expectedNumber int
	}{
		{0, color.RGBA64{'a'<<8 + 'b', 'c'<<8 + 'b', 'c', 'd'}, image.Point{}, []byte{'a', 'b', 'c', 'b', 0, 'c', 0, 'd'}, 8},
		{1, color.RGBA64{R: 'a'}, image.Point{}, []byte{'a', 0, 0, 0, 0, 0, 0}, 7},
		{2, color.RGBA64{R: 'a'}, image.Point{}, []byte{0, 0, 0, 0, 0, 0}, 6},
		{4, color.RGBA64{B: 'a', A: 'b'}, image.Point{}, []byte{0, 'a', 0, 'b'}, 4},
		{8, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, []byte{}, 0},
	}

	for i, test := range tests {
		b, n := SimplePoint64ReadWriter{}.Read(test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedBuff, b, "Test index %d", i)
	}
}

func Test_SimplePoint64ReadWriter_Write(t *testing.T) {
	tests := []struct {
		buff    []byte
		startOn int
		color   color.RGBA64
		point   image.Point

		expectedColor  *color.RGBA64
		expectedNumber int
	}{
		{[]byte{'a', 'b', 'c', 'd'}, 0, color.RGBA64{}, image.Point{}, &color.RGBA64{R: 'a'>>8 + 'b', G: 'c'>>8 + 'd'}, 4},
		{[]byte{'a', 'b', 'c', 'd'}, 1, color.RGBA64{}, image.Point{}, &color.RGBA64{R: 'b', G: 'c'>>8 + 'd'}, 3},
		{[]byte{'a', 'b', 'c', 'd'}, 2, color.RGBA64{R: 'f'}, image.Point{}, &color.RGBA64{'f', 'a'>>8 + 'b', 'c'>>8 + 'd', 0}, 4},
		{[]byte{'a', 'b', 'c', 'd'}, 6, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA64{'e', 'f', 'g', 'a'>>8 + 'b'}, 2},
		{[]byte{'a', 'b', 'c', 'd'}, 8, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA64{'e', 'f', 'g', 'h'}, 0},
		{[]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}, 0, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{},
			&color.RGBA64{'a'>>8 + 'b', 'c'>>8 + 'd', 'e'>>8 + 'f', 'g'>>8 + 'h'}, 8},
		{[]byte{'a'}, 0, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA64{'a'>>8 + 'e', 'f', 'g', 'h'}, 1},
		{[]byte{0, 'a'}, 0, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA64{'a', 'f', 'g', 'h'}, 2},
		{[]byte{'i'}, 2, color.RGBA64{'e', 'f', 'g', 'h'}, image.Point{}, &color.RGBA64{'e', 'i'>>8 + 'f', 'g', 'h'}, 1},
	}

	for i, test := range tests {
		c, n := SimplePoint64ReadWriter{}.Write(test.buff, test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedColor, c, "Test index %d", i)
	}
}
