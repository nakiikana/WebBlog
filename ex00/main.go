package main

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func createFace(size float64) font.Face {
	fnt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(fnt, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	return face
}

func drawString(img *image.RGBA, s string, x, y int, col color.Color) {

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: createFace(12),
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(x * 64),
			Y: fixed.Int26_6(y * 64),
		},
	}
	d.DrawString(s)
}

func main() {
	width := 300
	height := 300

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	purpleBeauty := color.RGBA{225, 216, 235, 255}
	sexyGreen := color.RGBA{157, 176, 157, 255}
	theColorOfMyLastGoodbye := color.RGBA{89, 99, 89, 255}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, purpleBeauty)
		}
	}
	x0 := 150
	y0 := 150
	r := 100
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)
	for x > y {
		for i := 0; i < 10; i++ {
			img.Set(x0+x+20, y0+y+20, sexyGreen)
			img.Set(x0+y+10, y0+x+10, sexyGreen)
			img.Set(x0-y-i, y0+x+i, sexyGreen)
			img.Set(x0-x-i, y0+y+i, sexyGreen)
			img.Set(x0-x-i, y0-y-i, sexyGreen)
			img.Set(x0-y-i, y0-x-i, sexyGreen)
			img.Set(x0+y+i, y0-x-i, sexyGreen)
			img.Set(x0+x+i, y0-y-i, sexyGreen)
		}

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
	drawString(img, "there are seeds in your soda..", 67, 150, theColorOfMyLastGoodbye)
	f, _ := os.Create("amazing_logo.png")
	png.Encode(f, img)
}
