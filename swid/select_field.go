package swid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// SelectField defines a select field widget.
type SelectField struct {
	widget.Select

	hovered        bool
	onHoverChanged func(bool)

	focused        bool
	onFocusChanged func(bool)
}

// NewSelectField creates a new select field widget.
func NewSelectField(options []string, changed func(string)) *SelectField {
	s := &SelectField{}
	s.ExtendBaseWidget(s)
	s.Options = options
	s.OnChanged = changed
	return s
}

// ===============================================================
// Implementation
// ===============================================================

// MinSize implements fyne.CanvasObject.
func (s *SelectField) MinSize() fyne.Size {
	s.ExtendBaseWidget(s)
	return s.BaseWidget.MinSize()
}

// FocusGained overrides widget.Select method.
func (s *SelectField) FocusGained() {
	s.focused = true
	s.Select.FocusGained()
	if s.onFocusChanged != nil {
		s.onFocusChanged(true)
	}
}

// FocusLost overrides widget.Select method.
func (s *SelectField) FocusLost() {
	s.focused = false
	s.Select.FocusLost()
	if s.onFocusChanged != nil {
		s.onFocusChanged(false)
	}
}

// MouseIn overrides widget.Select method.
func (s *SelectField) MouseIn(ev *desktop.MouseEvent) {
	s.hovered = true
	s.Select.MouseIn(ev)
	if s.onHoverChanged != nil {
		s.onHoverChanged(true)
	}
}

// MouseOut overrides widget.Select method.
func (s *SelectField) MouseOut() {
	s.hovered = false
	s.Select.MouseOut()
	if s.onHoverChanged != nil {
		s.onHoverChanged(false)
	}
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer implements fyne.Widget.
func (s *SelectField) CreateRenderer() fyne.WidgetRenderer {
	s.ExtendBaseWidget(s)
	return &selectFieldRenderer{
		WidgetRenderer: s.Select.CreateRenderer(),
		widget:         s,
	}
}

type selectFieldRenderer struct {
	fyne.WidgetRenderer
	widget *SelectField
}

func (r *selectFieldRenderer) Refresh() {
	r.WidgetRenderer.Refresh()
	bg := r.WidgetRenderer.Objects()[0].(*canvas.Rectangle)
	if r.widget.focused {
		bg.FillColor = theme.InputBackgroundColor()
		bg.Refresh()
	}
	if r.widget.hovered {
		bg.FillColor = theme.HoverColor()
		bg.Refresh()
	}

}
