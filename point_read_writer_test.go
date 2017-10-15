package imgio

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimplePointReadWriter_Read(t *testing.T) {
	b, n := SimplePointReadWriter{}.Read(&color.RGBA{'a', 'b', 'c', 'd'}, image.Point{})
	require.Equal(t, 4, n)
	require.Equal(t, []byte{'a', 'b', 'c', 'd'}, b)
}

func Test_SimplePointReadWriter_Write(t *testing.T) {
	c, n, err := SimplePointReadWriter{}.Write([]byte{'a', 'b', 'c', 'd'}, color.RGBA{}, image.Point{})
	require.Nil(t, err)
	require.Equal(t, 4, n)
	require.Equal(t, &color.RGBA{'a', 'b', 'c', 'd'}, c)
}
