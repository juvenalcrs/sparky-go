package swid

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TextFormField defines a special text field for Forms.
type TextFormField struct {
	widget.DisableableWidget
	BaseFormField

	Label       string
	TextStyle   fyne.TextStyle
	PlaceHolder string
	Hint        string
	Password    bool
	Wrapping    fyne.TextWrap
	Validator   fyne.StringValidator
	// ActionItem is a small item which is displayed at the outer right of the entry (like a password revealer)
	ActionItem fyne.CanvasObject
	MaxLength  int

	OnChanged func(s string)
	OnSaved   func(s string)

	labelAnim       *labelAnimation
	textField       *TextField
	initialText     string
	validationError error
	dirty           bool
}

// NewTextFormField creates a new special text field for Forms.
func NewTextFormField(label, initialText string) *TextFormField {
	t := &TextFormField{}
	t.ExtendBaseWidget(t)
	t.Label = label
	t.Wrapping = fyne.TextTruncate
	t.initialText = initialText
	t.setupTextField()
	return t
}

// Text returns the current text value.
func (t *TextFormField) Text() string {
	return t.textField.Text
}

// SetText manually sets the text of the TextFormField to the given text value.
func (t *TextFormField) SetText(text string) {
	t.textField.Text = text
	t.Refresh()
	if t.labelAnim != nil && !t.textField.focused && t.textField.Text == "" {
		t.labelAnim.Reverse()
	}
}

// Reset resets the text value to the initial value.
func (t *TextFormField) Reset() {
	t.dirty = false
	t.SetText(t.initialText)
	t.reset()
}

// Save triggers the OnSaved callback.
func (t *TextFormField) Save() {
	if t.OnSaved != nil {
		t.OnSaved(t.textField.Text)
	}
}

// Validate validates the field.
func (t *TextFormField) Validate() error {
	if t.Validator != nil {
		// means that this was called before CreateRenderer and
		// then Validator field is not copy to the textField yet,
		// so Refresh to generate it
		if t.textField.Validator == nil {
			t.ExtendBaseWidget(t)
			t.Refresh()
		}
		return t.validationError
	}
	return nil
}

func (t *TextFormField) setupTextField() {
	t.textField = NewTextField()
	t.textField.Text = t.initialText
	t.textField.OnChanged = func(s string) {
		if t.OnChanged != nil {
			t.OnChanged(s)
		}
		t.didChange()
	}
	t.textField.onFocusChanged = func(focused bool) {
		if t.labelAnim != nil && focused {
			t.labelAnim.Forward()
		}
		if t.labelAnim != nil && t.textField.Text == "" && !focused {
			t.labelAnim.Reverse()
		}
		if focused {
			t.dirty = true
		}
		t.Refresh()
	}
	t.textField.SetOnValidationChanged(func(e error) {
		t.validationError = e
		t.Refresh()
	})
}

// ===============================================================
// Renderer
// ===============================================================

// CreateRenderer implements fyne.Widget.
func (t *TextFormField) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)

	labelBg := canvas.NewRectangle(theme.InputBackgroundColor())
	label := canvas.NewText(t.Label, theme.PlaceHolderColor())
	// label.TextSize = theme.CaptionTextSize() // TODO
	hint := canvas.NewText(t.Hint, theme.PlaceHolderColor())
	hint.TextSize = hintTextSize()
	t.textField.Validator = t.Validator
	t.textField.Validate() // validates as soon as it is created

	r := &textFormFieldRenderer{
		labelBg: labelBg,
		label:   label,
		hint:    hint,
		widget:  t,
		objects: []fyne.CanvasObject{labelBg, label, t.textField, hint},
	}
	return r
}

type textFormFieldRenderer struct {
	labelBg *canvas.Rectangle
	label   *canvas.Text
	hint    *canvas.Text

	widget  *TextFormField
	objects []fyne.CanvasObject
}

func (r *textFormFieldRenderer) Destroy() {
	if r.widget.labelAnim != nil {
		r.widget.labelAnim.Stop()
	}
}

func (r *textFormFieldRenderer) Layout(size fyne.Size) {
	insetPad := fieldInsetPad()
	stackedLabelTextSize, _ := stackedLabelProps()
	stackedlabelMinHeight := fyne.MeasureText(r.label.Text, stackedLabelTextSize, r.label.TextStyle).Height
	r.labelBg.Move(fyne.NewPos(0, 0))
	r.labelBg.Resize(fyne.NewSize(size.Width, stackedlabelMinHeight-theme.InputBorderSize()))

	// If label animation is nil, it means we are in initial state, so setup
	if r.widget.labelAnim == nil {
		r.widget.labelAnim = newLabelAnimation(r.label, r.widget)
		labelPosY := float32(0)
		if r.widget.textField.Text != "" {
			r.label.TextSize, labelPosY = stackedLabelProps()
			r.label.Move(fyne.NewPos(insetPad, labelPosY))
		} else {
			r.label.TextSize, labelPosY = nonStackedLabelProps()
			r.label.Move(fyne.NewPos(insetPad, labelPosY))
		}
	}
	// r.label.Move(fyne.NewPos(insetPad, theme.InputBorderSize()*2)) // TODO
	// Use the label.MinSize() to use the current text size.
	r.label.Resize(fyne.NewSize(size.Width-2*insetPad, r.label.MinSize().Height))

	ypos := stackedlabelMinHeight - theme.InputBorderSize()*2
	textMinHeight := r.widget.textField.MinSize().Height
	r.widget.textField.Move(fyne.NewPos(0, ypos))
	r.widget.textField.Resize(fyne.NewSize(size.Width, textMinHeight))

	ypos += textMinHeight
	r.hint.Move(fyne.NewPos(insetPad, ypos))
	r.hint.Resize(fyne.NewSize(size.Width-2*insetPad, r.hint.MinSize().Height))
}

func (r *textFormFieldRenderer) MinSize() fyne.Size {
	min := r.widget.textField.MinSize()
	stackedLabelTextSize, _ := stackedLabelProps()
	labelMin := fyne.MeasureText(r.label.Text, stackedLabelTextSize, r.label.TextStyle)
	hintMin := r.hint.MinSize()
	min.Height += labelMin.Height - theme.InputBorderSize()*2
	min.Height += hintMin.Height
	min.Width = fyne.Max(min.Width, theme.Padding()*4+labelMin.Width)
	min.Width = fyne.Max(min.Width, theme.Padding()*4+hintMin.Width)
	return min
}

func (r *textFormFieldRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *textFormFieldRenderer) Refresh() {
	r.copyToTextField()

	r.labelBg.FillColor = theme.InputBackgroundColor()
	r.labelBg.Refresh()

	r.label.Text = r.widget.Label
	// r.label.TextSize = theme.CaptionTextSize() // TODO remove
	if r.widget.textField.focused {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.PlaceHolderColor()
	}

	r.hint.TextSize = hintTextSize()
	if !r.widget.textField.focused && r.widget.dirty && r.widget.validationError != nil {
		r.hint.Text = r.widget.validationError.Error()
		r.hint.Color = theme.ErrorColor()
		r.label.Color = theme.ErrorColor()
	} else {
		r.hint.Text = r.widget.Hint
		r.hint.Color = theme.PlaceHolderColor()
	}
	r.label.Refresh()
	r.hint.Refresh()
}

func (r *textFormFieldRenderer) copyToTextField() {
	r.widget.textField.TextStyle = r.widget.TextStyle
	r.widget.textField.ActionItem = r.widget.ActionItem
	r.widget.textField.Wrapping = r.widget.Wrapping
	r.widget.textField.Password = r.widget.Password
	r.widget.textField.PlaceHolder = r.widget.PlaceHolder
	r.widget.textField.Wrapping = r.widget.Wrapping
	r.widget.textField.MaxLength = r.widget.MaxLength
	r.widget.textField.Validator = r.widget.Validator
	if r.widget.Disabled() {
		r.widget.textField.Disable()
	} else {
		r.widget.textField.Enable()
	}
	r.widget.textField.Refresh()
}

// InsetPad for Label and Hint text inside the field
func fieldInsetPad() float32 {
	return 2*theme.Padding() - 1
}

func hintTextSize() float32 {
	return theme.CaptionTextSize() - 1
}

func stackedLabelProps() (textSize float32, posY float32) {
	return theme.CaptionTextSize(), theme.InputBorderSize() * 2
}

func nonStackedLabelProps() (textSize float32, posY float32) {
	return theme.TextSize(), theme.InputBorderSize() * 7
}

type labelAnimation struct {
	anim  *fyne.Animation
	label *canvas.Text
	w     fyne.Widget
}

func newLabelAnimation(label *canvas.Text, w fyne.Widget) *labelAnimation {
	return &labelAnimation{
		anim:  &fyne.Animation{Duration: canvas.DurationShort},
		label: label,
		w:     w,
	}
}

func (a *labelAnimation) animate(reverse bool) {
	startTextSize, startPosY := nonStackedLabelProps()
	endTextSize, endPosY := stackedLabelProps()
	deltaTextSize := endTextSize - startTextSize
	deltaPosY := endPosY - startPosY
	if reverse {
		startTextSize, endTextSize = endTextSize, startTextSize
		startPosY, endPosY = endPosY, startPosY
		deltaTextSize = -deltaTextSize
		deltaPosY = -deltaPosY
	}
	if a.label.Position().Y == endPosY {
		// return because it is already in the final position.
		return
	}
	insetPad := fieldInsetPad()
	a.label.Move(fyne.NewPos(insetPad, startPosY))
	a.anim.Tick = func(v float32) {
		a.label.TextSize = startTextSize + deltaTextSize*v
		a.label.Move(fyne.NewPos(insetPad, startPosY+deltaPosY*v))
		a.label.Refresh()
	}
	a.anim.Curve = fyne.AnimationEaseOut
	a.anim.Start()
}

func (a *labelAnimation) Forward() {
	a.anim.Stop()
	a.animate(false)
}

func (a *labelAnimation) Reverse() {
	a.anim.Stop()
	a.animate(true)
}

func (a *labelAnimation) Stop() {
	a.anim.Stop()
}
