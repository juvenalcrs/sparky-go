package slayout

import (
	"fyne.io/fyne/v2"
)

type minLayout struct {
	minWidth, minHeight float32
}

// NewMinWidthLayout creates a new layout with a minimum width.
func NewMinWidthLayout(minWidth float32) fyne.Layout {
	return &minLayout{minWidth: minWidth}
}

// NewMinHeightLayout creates a new layout with a minimum height.
func NewMinHeightLayout(minHeight float32) fyne.Layout {
	return &minLayout{minHeight: minHeight}
}

// Layout is called to pack all child objects into a specified size.
func (l *minLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topLeft := fyne.NewPos(0, 0)
	for _, child := range objects {
		child.Move(topLeft)
		child.Resize(size)
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
func (l *minLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	min := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		min = min.Max(child.MinSize())
	}
	min = min.Max(fyne.NewSize(l.minWidth, l.minHeight))
	return min
}
