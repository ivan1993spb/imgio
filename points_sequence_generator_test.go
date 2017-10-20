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

func Test_SimplePointsSequenceGenerator_Next_IncrementsCursor(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rectangle{},
		cursor: 0,
	}

	for i := uint64(0); i < 10; i++ {
		require.Equal(t, i, g.cursor)
		g.Next()
	}
}

func Test_SimplePointsSequenceGenerator_Rewind_ResetsCursor(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rectangle{},
		cursor: 10,
	}

	g.Rewind()
	require.Equal(t, uint64(0), g.cursor)
}

func Test_SimplePointsSequenceGenerator_Current_ReturnsValidCoordinatesByIndex(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rect(-5, -5, 1, 1),
		cursor: 0,
	}

	tests := map[uint64]image.Point{
		0: {-5, -5},
		4: {-1, -5},
		6: {-5, -4},
		9: {-2, -4},
	}

	for cursor, point := range tests {
		g.cursor = cursor
		require.Equal(t, point, g.Current())
	}
}

func Test_SimplePointsSequenceGenerator_Current_ReturnsValidCoordinatesByIndex_NoSquare(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		rect:   image.Rect(0, 0, 30, 100),
		cursor: 0,
	}

	tests := map[uint64]image.Point{
		0:  {0, 0},
		4:  {4, 0},
		29: {29, 0},
		30: {0, 1},
		31: {1, 1},
	}

	for cursor, point := range tests {
		g.cursor = cursor
		require.Equal(t, point, g.Current(), "Error: cursor=%d", cursor)
	}
}

func Test_SimplePointsSequenceGenerator_Seek(t *testing.T) {
	g := &SimplePointsSequenceGenerator{
		cursor: 0,
	}

	tests := []uint64{2, 3, 4, 2, 3, 4, 7, 2}

	for _, offset := range tests {
		g.Seek(offset)
		require.Equal(t, offset, g.cursor)
	}
}
