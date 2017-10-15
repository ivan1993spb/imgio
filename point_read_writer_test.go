package imgio

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimplePointReadWriter_Read(t *testing.T) {
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

	for _, test := range tests {
		b, n := SimplePointReadWriter{}.Read(test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n)
		require.Equal(t, test.expectedBuff, b)
	}
}

func Test_SimplePointReadWriter_Write(t *testing.T) {
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
	}

	for _, test := range tests {
		c, n := SimplePointReadWriter{}.Write(test.buff, test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n)
		require.Equal(t, test.expectedColor, c)
	}
}
