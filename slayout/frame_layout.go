package slayout

import (
	"fyne.io/fyne/v2"

	"github.com/fpabl0/sparky-go/internal/themedwid"
)

type frameLayout struct {
	border          *themedwid.ThemedBorder
	margin, padding float32
}

// NewFrameLayout creates a new frame layout.
func NewFrameLayout(border *themedwid.ThemedBorder, margin, padding float32) fyne.Layout {
	return &frameLayout{border, margin, padding}
}

func (l *frameLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	l.border.Move(fyne.NewPos(l.margin, l.margin))
	l.border.Resize(fyne.NewSize(size.Width-2*l.margin, size.Height-2*l.margin))

	insetPad := l.margin + l.border.StrokeWidth() + l.padding
	contentPos := fyne.NewPos(insetPad, insetPad)
	contentSize := fyne.NewSize(size.Width-2*insetPad, size.Height-2*insetPad)
	for _, child := range objects {
		if l.border == child {
			continue
		}
		child.Move(contentPos)
		child.Resize(contentSize)
	}
}

func (l *frameLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	contentMin := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() || l.border == child {
			continue
		}

		contentMin = contentMin.Max(child.MinSize())
	}

	insetPad := l.margin + l.border.StrokeWidth() + l.padding
	contentMin = contentMin.Add(fyne.NewSize(insetPad, insetPad))
	return contentMin
}
