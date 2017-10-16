package imgio

import (
	"image"
	"image/color"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Image_Write(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePointReadWriter{},
	}

	n, err := img.Write([]byte("testing"))
	require.Equal(t, 7, n)
	require.Nil(t, err)

	r1, g1, b1, a1 := img.img.At(0, 0).RGBA()
	r2, g2, b2, a2 := img.img.At(1, 0).RGBA()

	require.Equal(t, byte('t'), byte(r1))
	require.Equal(t, byte('e'), byte(g1))
	require.Equal(t, byte('s'), byte(b1))
	require.Equal(t, byte('t'), byte(a1))
	require.Equal(t, byte('i'), byte(r2))
	require.Equal(t, byte('n'), byte(g2))
	require.Equal(t, byte('g'), byte(b2))
	require.Equal(t, byte(0), byte(a2))
}

func Test_Image_Read(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePointReadWriter{},
	}

	img.img.Set(0, 0, &color.RGBA{0, 't', 'e', 's'})
	img.img.Set(1, 0, &color.RGBA{0, 't', 0, 'i'})
	img.img.Set(2, 0, &color.RGBA{'n', 'g', 0, 0})

	b, err := ioutil.ReadAll(img)
	b = b[:12]
	require.Nil(t, err, "Cannot read all from img")
	require.Equal(t, []byte{0, 't', 'e', 's', 0, 't', 0, 'i', 'n', 'g', 0, 0}, b)
}
