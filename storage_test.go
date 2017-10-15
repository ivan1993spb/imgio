package imgio

import (
	"image"
	"image/jpeg"
	"io"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	s := &Storage{
		img: image.NewRGBA(image.Rect(0, 0, 5, 5)),
		gen: NewSimplePointsSequenceGenerator(image.Rect(0, 0, 5, 5)),
		prw: SimplePointReadWriter{},
	}

	s.Write([]byte(`test123`))
	f, _ := os.Create("/tmp/image.jpeg")
	defer f.Close()

	io.Copy(os.Stdout, s)

	jpeg.Encode(f, s.img, &jpeg.Options{jpeg.DefaultQuality})
}
