package sparky

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

func infoColor() color.Color {
	return &color.NRGBA{R: 61, G: 194, B: 255, A: 255}
}

func successColor() color.Color {
	set := fyne.CurrentApp().Settings()
	lightBg := set.Theme().Color(theme.ColorNameBackground, theme.VariantLight)
	if theme.BackgroundColor() == lightBg {
		return &color.NRGBA{R: 45, G: 160, B: 111, A: 255}
	}
	return &color.NRGBA{R: 47, G: 223, B: 117, A: 255}
}

func dialogBackgroundColor() color.Color {
	rr, gg, bb, _ := theme.BackgroundColor().RGBA()
	bgColor := &color.NRGBA{R: uint8(rr), G: uint8(gg), B: uint8(bb), A: 230}
	return bgColor
}

func dialogTitleSize() float32 {
	return theme.TextSize() + 4
}

func dialogInsetPad() float32 {
	return 3 * theme.Padding()
}
