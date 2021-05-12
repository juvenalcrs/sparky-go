package scont

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/fpabl0/sparky-go/slayout"
)

// NewPadded creates a new padded container with top, bottom, left, right parameters.
func NewPadded(top, bottom, left, right float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(top, bottom, left, right), objects...)
}

// NewPaddedAll creates a new padded container with the same padding in all the sides.
func NewPaddedAll(all float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(all, all, all, all), objects...)
}

// NewPaddedTop creates a new top padded container.
func NewPaddedTop(top float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(top, 0, 0, 0), objects...)
}

// NewPaddedBottom creates a new bottom padded container.
func NewPaddedBottom(bottom float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(0, bottom, 0, 0), objects...)
}

// NewPaddedLeft creates a new left padded container.
func NewPaddedLeft(left float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(0, 0, left, 0), objects...)
}

// NewPaddedRight creates a new right padded container.
func NewPaddedRight(right float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(0, 0, 0, right), objects...)
}

// NewPaddedSym creates a new symmetric padded container.
func NewPaddedSym(vertical float32, horizontal float32, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(slayout.NewPaddedLayout(vertical, vertical, horizontal, horizontal), objects...)
}
