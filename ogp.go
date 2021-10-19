package ogpapp

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func createImage(width, height int, fontsize float64, ft *truetype.Font, text, out string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	opt := truetype.Options{
		Size: fontsize,
	}
	face := truetype.NewFace(ft, &opt)
	dr := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}
	x := (fixed.I(width) - dr.MeasureString(text)) / 2
	dr.Dot.X = x
	y := (height + int(fontsize)/2) / 2
	dr.Dot.Y = fixed.I(y)

	dr.DrawString(text)

	outfile, err := os.Create(out)
	if err != nil {
		return err
	}
	defer outfile.Close()

	if err := png.Encode(outfile, img); err != nil {
		return err
	}
	return nil
}
