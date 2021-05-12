package scont

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/fpabl0/sparky-go/internal/themedwid"
	"github.com/fpabl0/sparky-go/slayout"
)

// NewFrame creates a new frame container.
func NewFrame(margin, padding float32, objects ...fyne.CanvasObject) *fyne.Container {
	border := themedwid.NewThemedBorder(2.5)
	return container.New(slayout.NewFrameLayout(border, margin, padding), append(objects, border)...)
}
