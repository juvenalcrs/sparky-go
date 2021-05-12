package slayout

import (
	"fyne.io/fyne/v2"
)

type paddedLayout struct {
	top, bottom, left, right float32
}

// NewPaddedLayout creates a new padded layout.
func NewPaddedLayout(top, bottom, left, right float32) fyne.Layout {
	return &paddedLayout{top, bottom, left, right}
}

func (l *paddedLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	pos := fyne.NewPos(l.left, l.top)
	siz := fyne.NewSize(size.Width-l.left-l.right, size.Height-l.top-l.bottom)
	for _, child := range objects {
		child.Move(pos)
		child.Resize(siz)
	}
}

func (l *paddedLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	min := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		min = min.Max(child.MinSize())
	}
	min = min.Add(fyne.NewSize(l.left+l.right, l.top+l.bottom))
	return min
}
