package imgio

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"image"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ImageGroup_ReadWriteHash_OneImage(t *testing.T) {
	group := &ImageGroup{
		images: []*Image{
			{
				img: image.NewRGBA(image.Rect(0, 0, 100, 100)),
				gen: &SimplePointsSequenceGenerator{
					rect:   image.Rect(0, 0, 100, 100),
					cursor: 0,
				},
				prw: GentlePoint16ReadWriter{},
			},
		},
	}
	hasher := md5.New()
	buff := bytes.NewBuffer(nil)
	n, err := buff.ReadFrom(io.TeeReader(io.LimitReader(rand.Reader, group.Size()), hasher))
	require.Nil(t, err)
	require.Equal(t, group.Size(), n)
	firstSum := hasher.Sum(nil)
	hasher.Reset()
	n, err = buff.WriteTo(group)
	require.Nil(t, err)
	require.Equal(t, group.Size(), n)
	group.Rewind()
	n, err = io.Copy(hasher, group)
	require.Nil(t, err)
	require.Equal(t, group.Size(), n)
	secondSum := hasher.Sum(nil)
	require.Equal(t, firstSum, secondSum)
}

func Test_ImageGroup_ReadWriteHash_ManyImage(t *testing.T) {
	group := &ImageGroup{
		images: []*Image{
			{
				img: image.NewRGBA(image.Rect(0, 0, 100, 100)),
				gen: &SimplePointsSequenceGenerator{
					rect:   image.Rect(0, 0, 100, 100),
					cursor: 0,
				},
				prw: GentlePoint16ReadWriter{},
			},
			{
				img: image.NewRGBA(image.Rect(0, 0, 100, 100)),
				gen: &SimplePointsSequenceGenerator{
					rect:   image.Rect(0, 0, 100, 100),
					cursor: 0,
				},
				prw: SimplePoint32ReadWriter{},
			},
			{
				img: image.NewRGBA64(image.Rect(0, 0, 100, 100)),
				gen: &SimplePointsSequenceGenerator{
					rect:   image.Rect(0, 0, 100, 100),
					cursor: 0,
				},
				prw: SimplePoint64ReadWriter{},
			},
			{
				img: image.NewRGBA(image.Rect(0, 0, 100, 10)),
				gen: &SimplePointsSequenceGenerator{
					rect:   image.Rect(0, 0, 100, 10),
					cursor: 0,
				},
				prw: SimplePoint32ReadWriter{},
			},
		},
	}
	hasher := md5.New()
	buff := bytes.NewBuffer(nil)
	n, err := buff.ReadFrom(io.TeeReader(io.LimitReader(rand.Reader, group.Size()), hasher))
	require.Nil(t, err)
	require.Equal(t, group.Size(), int64(n))
	firstSum := hasher.Sum(nil)
	hasher.Reset()
	n, err = buff.WriteTo(group)
	require.Nil(t, err)
	require.Equal(t, group.Size(), n)
	group.Rewind()
	n, err = io.Copy(hasher, group)
	require.Nil(t, err)
	require.Equal(t, group.Size(), n)
	secondSum := hasher.Sum(nil)
	require.Equal(t, firstSum, secondSum)
}
