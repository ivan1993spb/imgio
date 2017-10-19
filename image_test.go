package imgio

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"image"
	"image/color"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Image_Write_UsePoint32(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
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

func Test_Image_Write_UsePoint32HandleErrOverflowOnePoint(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 1, 1)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 1, 1),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
	}

	n, err := img.Write([]byte("testing"))
	require.Equal(t, 4, n)
	require.Equal(t, ErrOverflow, err)
}

func Test_Image_Write_UsePoint32HandleErrOverflowManyPoints(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 10, 10)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 10, 10),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
	}

	size := img.Size() * 2
	buff := make([]byte, size)
	n, err := rand.Reader.Read(buff)
	require.Equal(t, size, int64(n))
	require.Nil(t, err)

	n, err = img.Write(buff)
	require.Equal(t, img.Size(), int64(n))
	require.Equal(t, ErrOverflow, err)
}

func Test_Image_Write_UsePoint64(t *testing.T) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}

	n, err := img.Write([]byte("testing"))
	require.Equal(t, 7, n)
	require.Nil(t, err)

	r1, g1, b1, a1 := img.img.At(0, 0).RGBA()
	r2, g2, b2, a2 := img.img.At(1, 0).RGBA()

	require.Equal(t, uint32('t'<<8+'e'), r1)
	require.Equal(t, uint32('s'<<8+'t'), g1)
	require.Equal(t, uint32('i'<<8+'n'), b1)
	require.Equal(t, uint32('g'<<8), a1)
	require.Equal(t, uint32(0), r2)
	require.Equal(t, uint32(0), g2)
	require.Equal(t, uint32(0), b2)
	require.Equal(t, uint32(0), a2)
}

func Test_Image_Write_UsePoint64_ErrOverflowOnePoint(t *testing.T) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(0, 0, 1, 1)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 1, 1),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}

	n, err := img.Write([]byte("testing 12345678"))
	require.Equal(t, 8, n)
	require.Equal(t, ErrOverflow, err)
}

func Test_Image_Write_UsePoint64HandleErrOverflowManyPoints(t *testing.T) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(0, 0, 10, 10)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 10, 10),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}

	size := img.Size() * 2
	buff := make([]byte, size)
	n, err := rand.Reader.Read(buff)
	require.Equal(t, size, int64(n))
	require.Nil(t, err)

	n, err = img.Write(buff)
	require.Equal(t, img.Size(), int64(n))
	require.Equal(t, ErrOverflow, err)
}

func Test_Image_Read_UsePoint32(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
	}

	img.img.Set(0, 0, &color.RGBA{0, 't', 'e', 's'})
	img.img.Set(1, 0, &color.RGBA{0, 't', 0, 'i'})
	img.img.Set(2, 0, &color.RGBA{'n', 'g', 0, 0})

	size := 12
	buff := make([]byte, size)
	n, err := img.Read(buff)
	require.Nil(t, err, "Cannot read from img")
	require.Equal(t, size, n)
	require.Equal(t, []byte{0, 't', 'e', 's', 0, 't', 0, 'i', 'n', 'g', 0, 0}, buff)
}

func Test_Image_Read_UsePoint64(t *testing.T) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(0, 0, 5, 5)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 5, 5),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}

	img.img.Set(0, 0, &color.RGBA64{0, 't' << 8, 'e', 's'})
	img.img.Set(1, 0, &color.RGBA64{0, 't' << 8, 0, 'i'})
	img.img.Set(2, 0, &color.RGBA64{'n', 'g', 'o'<<8 + 'k', 'e'<<8 + 'y'})

	size := 24
	buff := make([]byte, size)
	n, err := img.Read(buff)
	require.Nil(t, err, "Cannot read from img")
	require.Equal(t, size, n)
	require.Equal(t, []byte{
		0, 0, 't', 0, 0, 'e', 0, 's',
		0, 0, 't', 0, 0, 0, 0, 'i',
		0, 'n', 0, 'g', 'o', 'k', 'e', 'y',
	}, buff)
}

func Test_Image_ReadWrite32_Hash(t *testing.T) {
	img := &Image{
		img: image.NewRGBA(image.Rect(0, 0, 100, 100)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 100, 100),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
	}
	hasher := md5.New()
	buff := bytes.NewBuffer(nil)
	n, err := buff.ReadFrom(io.TeeReader(io.LimitReader(rand.Reader, img.Size()), hasher))
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	firstSum := hasher.Sum(nil)
	hasher.Reset()
	n, err = buff.WriteTo(img)
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	img.gen.Rewind()
	img.byteCursor = 0
	require.True(t, img.gen.Valid())
	n, err = io.Copy(hasher, img)
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	secondSum := hasher.Sum(nil)
	require.Equal(t, firstSum, secondSum)
}

func Test_Image_ReadWrite64_Hash(t *testing.T) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(0, 0, 100, 100)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(0, 0, 100, 100),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}
	hasher := md5.New()
	buff := bytes.NewBuffer(nil)
	n, err := buff.ReadFrom(io.TeeReader(io.LimitReader(rand.Reader, img.Size()), hasher))
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	firstSum := hasher.Sum(nil)
	hasher.Reset()
	n, err = buff.WriteTo(img)
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	img.gen.Rewind()
	img.byteCursor = 0
	require.True(t, img.gen.Valid())
	n, err = io.Copy(hasher, img)
	require.Nil(t, err)
	require.Equal(t, img.Size(), n)
	secondSum := hasher.Sum(nil)
	require.Equal(t, firstSum, secondSum)
}

func WriteBytesToImage64(x0, y0, x1, y1 int) (int64, error) {
	img := &Image{
		img: image.NewRGBA64(image.Rect(x0, y0, x1, y1)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(x0, y0, x1, y1),
			cursor: 0,
		},
		prw: SimplePoint64ReadWriter{},
	}
	return io.CopyN(img, rand.Reader, img.Size())
}

func Benchmark_WriteBytesToImage64(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := WriteBytesToImage64(0, 0, 1000, 1000)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func WriteBytesToImage32(x0, y0, x1, y1 int) (int64, error) {
	img := &Image{
		img: image.NewRGBA(image.Rect(x0, y0, x1, y1)),
		gen: &SimplePointsSequenceGenerator{
			rect:   image.Rect(x0, y0, x1, y1),
			cursor: 0,
		},
		prw: SimplePoint32ReadWriter{},
	}
	return io.CopyN(img, rand.Reader, img.Size())
}

func Benchmark_WriteBytesToImage32(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, err := WriteBytesToImage32(0, 0, 1000, 1000)
		if err != nil {
			b.Fatal(err)
		}
	}
}
