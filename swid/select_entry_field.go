package swid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// SelectEntryField defines a select entry field widget.
type SelectEntryField struct {
	widget.SelectEntry

	focused        bool
	onFocusChanged func(bool)
}

// NewSelectEntryField creates a new select entry field.
func NewSelectEntryField(options []string) *SelectEntryField {
	s := &SelectEntryField{}
	s.ExtendBaseWidget(s)
	s.Wrapping = fyne.TextTruncate
	s.SetOptions(options)
	return s
}

// ===============================================================
// Implementation
// ===============================================================

// MinSize implements fyne.CanvasObject.
func (s *SelectEntryField) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// FocusGained overrides widget.Select method.
func (s *SelectEntryField) FocusGained() {
	s.focused = true
	s.SelectEntry.FocusGained()
	if s.onFocusChanged != nil {
		s.onFocusChanged(true)
	}
}

// FocusLost overrides widget.Select method.
func (s *SelectEntryField) FocusLost() {
	s.focused = false
	s.SelectEntry.FocusLost()
	if s.onFocusChanged != nil {
		s.onFocusChanged(false)
	}
}
