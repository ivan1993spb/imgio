package imgio

import (
	"fmt"
	"image"
	"testing"
)

func Test_SimplePointsSequenceGenerator_Current(t *testing.T) {
	g := NewSimplePointsSequenceGenerator(image.Rect(-10, -10, -1, -1))
	for i := 0; i < g.rect.Size().X*g.rect.Size().Y; i++ {
		p := g.Current()
		fmt.Printf("%4d %s\n", i, p)
		g.Next()
	}
}
