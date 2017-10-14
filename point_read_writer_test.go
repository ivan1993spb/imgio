package imgio

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimplePointReadWriter_Read(t *testing.T) {
	b := SimplePointReadWriter{}.Read(&color.RGBA{'a', 'b', 'c', 'd'})
	require.Equal(t, []byte{'a', 'b', 'c', 'd'}, b)
}

func Test_SimplePointReadWriter_Write(t *testing.T) {
	c, err := SimplePointReadWriter{}.Write([]byte{'a', 'b', 'c', 'd'})
	require.Nil(t, err)
	require.Equal(t, &color.RGBA{'a', 'b', 'c', 'd'}, c)
}
