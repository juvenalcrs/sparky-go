package scont

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/fpabl0/sparky-go/slayout"
)

// NewMinWidth creates a container with a minimum width.
func NewMinWidth(minWidth float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewMinWidthLayout(minWidth), objects...)
}

// NewMinHeight creates a container with a minimum height.
func NewMinHeight(minHeight float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewMinHeightLayout(minHeight), objects...)
}
