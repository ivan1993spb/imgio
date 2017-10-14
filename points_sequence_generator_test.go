package imgio

import (
	"image"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimplePointsSequenceGenerator_Valid_ReturnsTrue(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rect(0, 0, 1, 1),
		cursor: 0,
	}
	require.True(t, g.Valid())
}

func Test_SimplePointsSequenceGenerator_Valid_RetursFalse(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rect(-1, -1, 1, 1),
		cursor: 4,
	}
	require.False(t, g.Valid())
}

func Test_SimplePointsSequenceGenerator_Current_ReturnsValidCoordinatesByIndex(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rect(-5, -5, 1, 1),
		cursor: 0,
	}

	tests := map[uint64]image.Point{
		0: image.Point{-5, -5},
		4: image.Point{-1, -5},
		6: image.Point{-5, -4},
		9: image.Point{-2, -4},
	}

	for cursor, point := range tests {
		g.cursor = cursor
		require.Equal(t, point, g.Current())
	}
}
