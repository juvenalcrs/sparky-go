package scont

import (
	"fyne.io/fyne/v2"
)

// LayoutBuilder defines layout constructor.
type LayoutBuilder func([]fyne.CanvasObject, fyne.Size)

// NewStack creates a stack container.
//
//  scont.NewStack(objects, func(objs []fyne.CanvasObject, size fyne.Size) {
//     label.Move(...)
//     label.Resize(...)
//  })
//
func NewStack(objects []fyne.CanvasObject, layoutBuilder LayoutBuilder) *fyne.Container {
	// TODO
	return nil
}
