package go_hugipipes_signal_drawer

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
)

type Drawable interface {
	Set(x, y int, c color.Color)
	DrawString(x, y int, text string, c color.Color)
}

type ImageDrawable struct {
	img *image.RGBA
}

func NewImageDrawable(img *image.RGBA) *ImageDrawable {
	return &ImageDrawable{
		img: img,
	}
}

func (s *ImageDrawable) DrawString(x, y int, text string, c color.Color) {
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	fd := &font.Drawer{
		Dst:  s.img,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	fd.DrawString(text)
}
func (s *ImageDrawable) Set(x, y int, c color.Color) {
	s.img.Set(x, y, c)
}
