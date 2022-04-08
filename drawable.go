package go_hugipipes_signal_drawer

import "image/color"

type Drawable interface {
	Set(x, y int, c color.Color)
	DrawString(x, y int, text string, c color.Color)
}
