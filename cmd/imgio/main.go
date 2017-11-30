package main

import (
	"bytes"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"sort"

	"github.com/ivan1993spb/imgio"

	"github.com/urfave/cli"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{}

	app.Commands = []cli.Command{
		{
			Name: "encode",
			Action: func(c *cli.Context) error {
				img := imgio.NewImage(
					image.NewRGBA(image.Rect(0, 0, 10, 10)),
					imgio.NewSimplePointsSequenceGenerator(image.Rect(0, 0, 10, 10)),
					imgio.SimplePoint32ReadWriter{},
				)
				n, err := io.Copy(img, os.Stdin)
				log.Println(n, err)

				err = jpeg.Encode(os.Stdout, img, &jpeg.Options{})
				log.Println(err)
				return nil
			},
		},
		{
			Name: "decode",
			Action: func(c *cli.Context) error {
				i, err := jpeg.Decode(os.Stdin)
				log.Println(err)

				cimg := image.NewRGBA(i.Bounds())
				draw.Draw(cimg, i.Bounds(), i, image.Point{}, draw.Over)

				img := imgio.NewImage(
					cimg,
					imgio.NewSimplePointsSequenceGenerator(image.Rect(0, 0, 10, 10)),
					imgio.SimplePoint32ReadWriter{},
				)
				n, err := io.Copy(os.Stdout, img)
				log.Println()
				log.Println(n, err)

				return nil
			},
		},

		{
			Name: "show",
			Action: func(c *cli.Context) error {
				buff := bytes.NewBuffer(nil)
				n, err := buff.ReadFrom(os.Stdin)
				log.Println("read", n, err)

				var img image.Image
				img, err = jpeg.Decode(buff)
				show(img)
				log.Printf("jpeg %#v %s\n", img, err)
				img, err = png.Decode(buff)
				show(img)
				log.Printf("png %#v %s\n", img, err)
				img, err = gif.Decode(buff)
				show(img)
				log.Printf("gif %#v %s\n", img, err)
				img, err = bmp.Decode(buff)
				show(img)
				log.Printf("bmp %#v %s\n", img, err)
				img, err = tiff.Decode(buff)
				show(img)
				log.Printf("tif %#v %s\nf", img, err)

				return nil
			},
		},

		{
			Flags: []cli.Flag{&cli.UintFlag{Name: "test "}},
			Name:  "encode_jpeg",
			Action: func(c *cli.Context) error {
				log.Println(c.Uint("test"))
				i := image.NewYCbCr(image.Rect(0, 0, 3, 3), image.YCbCrSubsampleRatio420)
				// i.YStride = 16
				// i.CStride = 8
				imgrw := imgio.NewImageReadWriterYCbCr(
					i,
					imgio.NewSimplePointsSequenceGenerator(image.Rect(0, 0, 3, 3)),
					imgio.PointReadWriterYCbCrSimple{uint8(c.Uint("test"))},
				)
				n, err := io.Copy(imgrw, os.Stdin)
				log.Println(n, err)
				log.Printf("%#v\n", i)

				err = jpeg.Encode(os.Stdout, i, &jpeg.Options{Quality: 100})
				log.Println(err)
				return nil
			},
		},
		{
			Name: "decode_jpeg",
			Action: func(c *cli.Context) error {
				img, err := jpeg.Decode(os.Stdin)
				log.Println(err)

				if imgYCbCr, ok := img.(*image.YCbCr); ok {
					imgrw := imgio.NewImageReadWriterYCbCr(
						imgYCbCr,
						imgio.NewSimplePointsSequenceGenerator(imgYCbCr.Rect),
						imgio.PointReadWriterYCbCrSimple{},
					)
					log.Printf("%#v\n", imgYCbCr)

					n, err := io.Copy(os.Stdout, imgrw)
					log.Println(n, err)

				} else {
					log.Println("not ok =(")
				}

				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func show(img image.Image) {
	switch img.(type) {
	case *image.Alpha:
		log.Println("Alpha")
	case *image.Alpha16:
		log.Println("Alpha16")
	case *image.CMYK:
		log.Println("CMYK")
	case *image.Gray:
		log.Println("Gray")
	case *image.Gray16:
		log.Println("Gray16")
	case *image.NRGBA:
		log.Println("NRGBA")
	case *image.NRGBA64:
		log.Println("NRGBA64")
	case *image.NYCbCrA:
		log.Println("NYCbCrA")
	case *image.Paletted:
		log.Println("Paletted")
	case *image.RGBA:
		log.Println("RGBA")
	case *image.RGBA64:
		log.Println("RGBA64")
	case *image.YCbCr:
		log.Println("YCbCr")
	default:
		log.Println("default")
	}
}
