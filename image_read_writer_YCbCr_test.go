package imgio

import (
	"crypto/rand"
	"image"
	"testing"

	"gopkg.in/stretchr/testify.v1/require"
)

func Test_ImageReadWriterYCbCr_Write_WriteUsingPointReadWriterYCbCrSimple(t *testing.T) {
	imgrw := &ImageReadWriterYCbCr{
		img: image.NewYCbCr(image.Rect(0, 0, 5, 2), image.YCbCrSubsampleRatio444),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 2),
			cursor: 0,
		},
		prw: PointReadWriterYCbCrSimple{},
	}

	n, err := imgrw.Write([]byte("testing"))
	require.Equal(t, 7, n)
	require.Nil(t, err)

	color01 := imgrw.img.YCbCrAt(0, 0)
	color02 := imgrw.img.YCbCrAt(1, 0)
	color03 := imgrw.img.YCbCrAt(2, 0)

	require.Equal(t, byte('t'), color01.Y)
	require.Equal(t, byte('e'), color01.Cb)
	require.Equal(t, byte('s'), color01.Cr)
	require.Equal(t, byte('t'), color02.Y)
	require.Equal(t, byte('i'), color02.Cb)
	require.Equal(t, byte('n'), color02.Cr)
	require.Equal(t, byte('g'), color03.Y)
	require.Equal(t, byte(0), color03.Cb)
	require.Equal(t, byte(0), color03.Cr)
}

func Test_ImageReadWriterYCbCr_Write_WriteUsingPointReadWriterYCbCrSimple_OnePoint_ExpectsOverflow(t *testing.T) {
	imgrw := &ImageReadWriterYCbCr{
		img: image.NewYCbCr(image.Rect(0, 0, 1, 1), image.YCbCrSubsampleRatio444),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 1, 1),
			cursor: 0,
		},
		prw: PointReadWriterYCbCrSimple{},
	}

	n, err := imgrw.Write([]byte("testing"))
	require.Equal(t, 3, n)
	require.Equal(t, ErrImageReadWriterYCbCrOverflow, err)
}

func Test_ImageReadWriterYCbCr_Write_WriteUsingPointReadWriterYCbCrSimple_ManyPoints_ExpectsOverflow(t *testing.T) {
	width := 10
	height := 5
	size := width * height * PointReadWriterYCbCrSimpleCapacity
	buffSize := size * 2

	imgrw := &ImageReadWriterYCbCr{
		img: image.NewYCbCr(image.Rect(0, 0, width, height), image.YCbCrSubsampleRatio444),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, width, height),
			cursor: 0,
		},
		prw: PointReadWriterYCbCrSimple{},
	}

	buff := make([]byte, buffSize)
	n, err := rand.Reader.Read(buff)
	require.EqualValues(t, buffSize, n)
	require.Nil(t, err)

	n, err = imgrw.Write(buff)
	require.EqualValues(t, size, n)
	require.Equal(t, ErrImageReadWriterYCbCrOverflow, err)
}
