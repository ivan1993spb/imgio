package imgio

import (
	"image"
	"image/color"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func Test_PointReadWriterYCbCrSimple_Read(t *testing.T) {
	tests := []struct {
		startOn int
		color   color.YCbCr
		point   image.Point

		expectedBuff   []byte
		expectedNumber int
	}{
		{0, color.YCbCr{0xff, 'b', 'c'}, image.Point{}, []byte{'b', 'c'}, 2},
		{1, color.YCbCr{Y: 'a'}, image.Point{}, []byte{0}, 1},
		{1, color.YCbCr{Cb: 'a', Cr: 'b'}, image.Point{}, []byte{'b'}, 1},
		{3, color.YCbCr{'e', 'f', 'g'}, image.Point{}, []byte{}, 0},
		{4, color.YCbCr{'e', 'f', 'g'}, image.Point{}, []byte{}, 0},
	}

	for i, test := range tests {
		b, n := PointReadWriterYCbCrSimple{}.Read(test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedBuff, b, "Test index %d", i)
	}
}

func Test_PointReadWriterYCbCrSimple_Write(t *testing.T) {
	tests := []struct {
		buff    []byte
		startOn int
		color   color.YCbCr
		point   image.Point

		expectedColor  color.YCbCr
		expectedNumber int
	}{
		{[]byte{'a', 'b', 'c', 'd'}, 0, color.YCbCr{}, image.Point{}, color.YCbCr{0xff, 'a', 'b'}, 2},
		{[]byte{'a', 'b', 'c', 'd'}, 1, color.YCbCr{Y: 'a'}, image.Point{}, color.YCbCr{0xff, 0, 'a'}, 1},
		{[]byte{'a', 'b', 'c', 'd'}, 2, color.YCbCr{Cb: 'a'}, image.Point{}, color.YCbCr{Cb: 'a'}, 0},
		{[]byte{'a', 'b', 'c', 'd'}, 4, color.YCbCr{'e', 'f', 'g'}, image.Point{}, color.YCbCr{'e', 'f', 'g'}, 0},
		{[]byte{'a', 'b', 'c', 'd', 'e', 'f'}, 0, color.YCbCr{'e', 'f', 'g'}, image.Point{}, color.YCbCr{0xff, 'b', 'c'}, 3},
		{[]byte{'a'}, 0, color.YCbCr{'e', 'f', 'g'}, image.Point{}, color.YCbCr{0xff, 'f', 'g'}, 1},
		{[]byte{'i'}, 2, color.YCbCr{'e', 'f', 'g'}, image.Point{}, color.YCbCr{0xff, 'f', 'i'}, 1},
		{[]byte{}, 0, color.YCbCr{'e', 'f', 'g'}, image.Point{}, color.YCbCr{0xff, 'f', 'g'}, 0},
	}

	for i, test := range tests {
		c, n := PointReadWriterYCbCrSimple{}.Write(test.buff, test.startOn, test.color, test.point)
		require.Equal(t, test.expectedNumber, n, "Test index %d", i)
		require.Equal(t, test.expectedColor, c, "Test index %d", i)
	}
}

func Test_PointReadWriterYCbCrSimple_Size(t *testing.T) {
	require.EqualValues(t, PointReadWriterYCbCrSimpleCapacity, PointReadWriterYCbCrSimple{}.Size(image.Point{}))
}
