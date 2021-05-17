package swid

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// FormField defines a widget that can be used inside a Form.
type FormField interface {
	fyne.Widget
	Reset()
	Save()
	Validate() error

	setParentForm(f *Form)
	didChange()
}

// BaseFormField defines a base form field.
type BaseFormField struct {
	widget.DisableableWidget
	Label string
	Hint  string

	labelAnim       *labelAnimation
	dirty           bool
	validationError error
	form            *Form

	impl fyne.Widget
}

// ExtendBaseFormField extends a base form field.
func (b *BaseFormField) ExtendBaseFormField(w fyne.Widget) {
	if b.impl != nil {
		return
	}
	b.ExtendBaseWidget(w)
	b.impl = w
}

func (b *BaseFormField) setParentForm(f *Form) {
	b.form = f
}

func (b *BaseFormField) didChange() {
	if b.form == nil {
		return
	}
	b.form.fieldDidChange()
}

// ===============================================================
// BaseRenderer
// ===============================================================

// CreateBaseRenderer creates a base form field renderer.
func (b *BaseFormField) CreateBaseRenderer(
	labelText, hintText string, fieldWidget fyne.Widget,
	isFieldEmpty func() bool,
	isFieldFocused func() bool,
	updateInternalField func(),
) fyne.WidgetRenderer {
	labelBg := newLabelBackground(theme.InputBackgroundColor(), fieldWidget)
	label := canvas.NewText(labelText, theme.PlaceHolderColor())
	hint := canvas.NewText(hintText, theme.PlaceHolderColor())
	hint.TextSize = hintTextSize()
	return &formFieldRenderer{
		labelBg:             labelBg,
		label:               label,
		fieldWidget:         fieldWidget,
		hint:                hint,
		labelBgColor:        func() color.Color { return theme.InputBackgroundColor() },
		isFieldEmpty:        isFieldEmpty,
		isFieldFocused:      isFieldFocused,
		updateInternalField: updateInternalField,
		formField:           b,
		objects:             []fyne.CanvasObject{labelBg, label, fieldWidget, hint},
	}
}

type formFieldRenderer struct {
	labelBg     *labelBackground
	label       *canvas.Text
	fieldWidget fyne.Widget
	hint        *canvas.Text

	labelBgColor        func() color.Color
	isFieldEmpty        func() bool
	isFieldFocused      func() bool
	updateInternalField func()

	formField *BaseFormField
	objects   []fyne.CanvasObject
}

func (r *formFieldRenderer) Destroy() {
	if r.formField.labelAnim != nil {
		r.formField.labelAnim.Stop()
	}
}

func (r *formFieldRenderer) Layout(size fyne.Size) {
	insetPad := r.fieldInsetPad()
	stackedLabelTextSize, _ := r.stackedLabelProps()
	stackedlabelMinHeight := fyne.MeasureText(r.label.Text, stackedLabelTextSize, r.label.TextStyle).Height
	r.labelBg.Move(fyne.NewPos(0, 0))
	r.labelBg.Resize(fyne.NewSize(size.Width, stackedlabelMinHeight-theme.InputBorderSize()))

	// If label animation is nil, it means we are in initial state, so setup
	if r.formField.labelAnim == nil {
		r.formField.labelAnim = r.newLabelAnimation()
		labelPosY := float32(0)
		if !r.isFieldEmpty() {
			r.label.TextSize, labelPosY = r.stackedLabelProps()
			r.label.Move(fyne.NewPos(insetPad, labelPosY))
		} else {
			r.label.TextSize, labelPosY = r.nonStackedLabelProps()
			r.label.Move(fyne.NewPos(insetPad, labelPosY))
		}
	}

	// Use the label.MinSize() to use the current text size.
	r.label.Resize(fyne.NewSize(size.Width-2*insetPad, r.label.MinSize().Height))

	ypos := stackedlabelMinHeight - theme.InputBorderSize()*2
	fieldMinHeight := r.fieldWidget.MinSize().Height
	r.fieldWidget.Move(fyne.NewPos(0, ypos))
	r.fieldWidget.Resize(fyne.NewSize(size.Width, fieldMinHeight))

	ypos += fieldMinHeight
	r.hint.Move(fyne.NewPos(insetPad, ypos))
	r.hint.Resize(fyne.NewSize(size.Width-2*insetPad, r.hint.MinSize().Height))
}

func (r *formFieldRenderer) MinSize() fyne.Size {
	min := r.fieldWidget.MinSize()
	stackedLabelTextSize, _ := r.stackedLabelProps()
	labelMin := fyne.MeasureText(r.label.Text, stackedLabelTextSize, r.label.TextStyle)
	hintMin := r.hint.MinSize()
	min.Height += labelMin.Height - theme.InputBorderSize()*2
	min.Height += hintMin.Height
	min.Width = fyne.Max(min.Width, theme.Padding()*4+labelMin.Width)
	min.Width = fyne.Max(min.Width, theme.Padding()*4+hintMin.Width)
	return min
}

func (r *formFieldRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *formFieldRenderer) Refresh() {
	if r.isFieldFocused() || !r.isFieldEmpty() {
		r.formField.dirty = true
	}
	if r.formField.labelAnim != nil && (r.isFieldFocused() || !r.isFieldEmpty()) {
		r.formField.labelAnim.Forward()
	}
	if r.formField.labelAnim != nil && r.isFieldEmpty() && !r.isFieldFocused() {
		r.formField.labelAnim.Reverse()
	}

	r.updateInternalField()

	r.labelBg.FillColor = r.labelBgColor()
	r.labelBg.Refresh()

	r.label.Text = r.formField.Label
	if r.isFieldFocused() {
		r.label.Color = theme.PrimaryColor()
	} else {
		r.label.Color = theme.PlaceHolderColor()
	}

	r.hint.TextSize = hintTextSize()
	if !r.isFieldFocused() && r.formField.dirty && r.formField.validationError != nil {
		r.hint.Text = r.formField.validationError.Error()
		r.hint.Color = theme.ErrorColor()
		r.label.Color = theme.ErrorColor()
	} else {
		r.hint.Text = r.formField.Hint
		r.hint.Color = theme.PlaceHolderColor()
	}
	r.label.Refresh()
	r.hint.Refresh()
}

// InsetPad for Label and Hint text inside the field
func (r *formFieldRenderer) fieldInsetPad() float32 {
	return 2 * theme.Padding()
}

func (r *formFieldRenderer) stackedLabelProps() (textSize float32, posY float32) {
	return theme.CaptionTextSize(), theme.InputBorderSize() * 2
}

func (r *formFieldRenderer) nonStackedLabelProps() (textSize float32, posY float32) {
	textSize = theme.TextSize()
	nonStackedMinHeight := fyne.MeasureText(r.label.Text, textSize, r.label.TextStyle).Height
	return textSize,
		(r.MinSize().Height - theme.InputBorderSize() - r.hint.MinSize().Height - nonStackedMinHeight) / 2
}

func hintTextSize() float32 {
	return theme.CaptionTextSize() - 1
}

// ===============================================================
// Label Background
// ===============================================================

type labelBackground struct {
	widget.BaseWidget
	FillColor color.Color

	hovered     bool
	hoverable   desktop.Hoverable
	cursor      desktop.Cursor
	fieldWidget fyne.Widget
}

func newLabelBackground(color color.Color, fieldWidget fyne.Widget) *labelBackground {
	b := &labelBackground{}
	b.ExtendBaseWidget(b)
	b.FillColor = color
	b.fieldWidget = fieldWidget
	b.cursor = desktop.DefaultCursor
	if cursorable, ok := fieldWidget.(desktop.Cursorable); ok {
		b.cursor = cursorable.Cursor()
	}
	if hoverable, ok := fieldWidget.(desktop.Hoverable); ok {
		b.hoverable = hoverable
	}
	return b
}

func (b *labelBackground) Cursor() desktop.Cursor {
	return b.cursor
}

func (b *labelBackground) MouseIn(ev *desktop.MouseEvent) {
	if b.hoverable == nil {
		return
	}
	b.hoverable.MouseIn(ev)
	b.hovered = true
	b.Refresh()
}

func (b *labelBackground) MouseMoved(*desktop.MouseEvent) {}

func (b *labelBackground) MouseOut() {
	if b.hoverable == nil {
		return
	}
	b.hoverable.MouseOut()
	b.hovered = false
	b.Refresh()
}

func (b *labelBackground) Tapped(ev *fyne.PointEvent) {
	if focusable, ok := b.fieldWidget.(fyne.Focusable); ok {
		cnv := fyne.CurrentApp().Driver().CanvasForObject(b)
		cnv.Focus(focusable)
	}
	if tappable, ok := b.fieldWidget.(fyne.Tappable); ok {
		tappable.Tapped(ev)
	}
}

func (b *labelBackground) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	rect := canvas.NewRectangle(b.FillColor)
	return &labelBackgroundRenderer{
		rect:    rect,
		widget:  b,
		objects: []fyne.CanvasObject{rect},
	}
}

type labelBackgroundRenderer struct {
	rect    *canvas.Rectangle
	widget  *labelBackground
	objects []fyne.CanvasObject
}

func (r *labelBackgroundRenderer) Destroy() {}

func (r *labelBackgroundRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *labelBackgroundRenderer) MinSize() fyne.Size {
	return r.rect.MinSize()
}

func (r *labelBackgroundRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *labelBackgroundRenderer) Refresh() {
	r.rect.FillColor = r.widget.FillColor
	if r.widget.hovered {
		r.rect.FillColor = theme.HoverColor()
	}
	r.rect.Refresh()
}

// ===============================================================
// Label animation
// ===============================================================

type labelAnimation struct {
	anim     *fyne.Animation
	renderer *formFieldRenderer
}

func (r *formFieldRenderer) newLabelAnimation() *labelAnimation {
	return &labelAnimation{
		anim:     &fyne.Animation{Duration: canvas.DurationShort},
		renderer: r,
	}
}

func (a *labelAnimation) animate(reverse bool) {
	startTextSize, startPosY := a.renderer.nonStackedLabelProps()
	endTextSize, endPosY := a.renderer.stackedLabelProps()
	deltaTextSize := endTextSize - startTextSize
	deltaPosY := endPosY - startPosY
	if reverse {
		startTextSize, endTextSize = endTextSize, startTextSize
		startPosY, endPosY = endPosY, startPosY
		deltaTextSize = -deltaTextSize
		deltaPosY = -deltaPosY
	}
	if a.renderer.label.Position().Y == endPosY {
		// return because it is already in the final position.
		return
	}
	insetPad := a.renderer.fieldInsetPad()
	a.renderer.label.Move(fyne.NewPos(insetPad, startPosY))
	a.anim.Tick = func(v float32) {
		a.renderer.label.TextSize = startTextSize + deltaTextSize*v
		a.renderer.label.Move(fyne.NewPos(insetPad, startPosY+deltaPosY*v))
		a.renderer.label.Refresh()
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
