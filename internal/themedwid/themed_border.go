package themedwid

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// ThemedBorder defines themed border widget.
type ThemedBorder struct {
	widget.BaseWidget
	strokeWidth float32
}

// NewThemedBorder creates a new themed border.
func NewThemedBorder(strokeWidth float32) *ThemedBorder {
	b := &ThemedBorder{}
	b.ExtendBaseWidget(b)
	b.strokeWidth = strokeWidth
	return b
}

// StrokeWidth returns the border stroke width.
func (b *ThemedBorder) StrokeWidth() float32 {
	return b.strokeWidth
}

// CreateRenderer implements fyne.Widget.
func (b *ThemedBorder) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = b.strokeWidth
	border.StrokeColor = theme.ErrorColor()
	r := &themedBorderRenderer{
		border:  border,
		objects: []fyne.CanvasObject{border},
	}
	r.Refresh()
	return r
}

type themedBorderRenderer struct {
	border  *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *themedBorderRenderer) Destroy() {}

func (r *themedBorderRenderer) Layout(size fyne.Size) {
	r.border.Resize(size)
}

func (r *themedBorderRenderer) MinSize() fyne.Size {
	return r.border.MinSize()
}

func (r *themedBorderRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *themedBorderRenderer) Refresh() {
	set := fyne.CurrentApp().Settings()
	lightBg := set.Theme().Color(theme.ColorNameBackground, theme.VariantLight)
	if theme.BackgroundColor() == lightBg {
		r.border.StrokeColor = &color.NRGBA{R: 215, G: 216, B: 218, A: 255}
	} else {
		r.border.StrokeColor = &color.NRGBA{R: 34, G: 36, B: 40, A: 255}
	}
	r.border.Refresh()
}
