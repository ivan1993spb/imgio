package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"os"
	"sort"

	"github.com/ivan1993spb/imgio"

	"github.com/urfave/cli"
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
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}
