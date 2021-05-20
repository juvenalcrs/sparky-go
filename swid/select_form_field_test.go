package swid

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestSelectFormField_Reset(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	sf := NewSelectFormField("Name", "Andrea", []string{"Peter", "Andrea"})
	sf.Placeholder = "Your name"

	assert.Equal(t, "Andrea", sf.Selected())

	sf.SetSelected("Peter")
	assert.Equal(t, "Peter", sf.Selected())

	sf.Reset()
	assert.Equal(t, "Andrea", sf.Selected())

	sf = NewSelectFormField("Name", "", []string{"Andrea", "Peter"})
	sf.Placeholder = "Your name"

	w := test.NewWindow(sf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	assert.Equal(t, "", sf.Selected())

	sf.SetSelected("Peter")
	assert.Equal(t, "Peter", sf.Selected())

	sf.Reset()
	assert.Equal(t, "", sf.Selected())

	w.Canvas().Focus(sf.selectField)
	assert.Equal(t, "", sf.Selected())
}

func TestSelectFormField_Hover(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	sf := NewSelectFormField("Name", "Andrea", []string{"Peter", "Andrea"})
	sf.Placeholder = "Your name"
	sf.Hint = "A hint text"

	w := test.NewWindow(sf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	sf.selectField.MouseIn(&desktop.MouseEvent{})
	test.AssertImageMatches(t, "select_form_field/hovered.png", w.Canvas().Capture())

	sf.selectField.MouseOut()
	test.AssertImageMatches(t, "select_form_field/unhovered.png", w.Canvas().Capture())

	r := test.WidgetRenderer(sf).(*formFieldRenderer)

	r.labelBg.MouseIn(&desktop.MouseEvent{})
	test.AssertImageMatches(t, "select_form_field/hovered.png", w.Canvas().Capture())

	r.labelBg.MouseOut()
	test.AssertImageMatches(t, "select_form_field/unhovered.png", w.Canvas().Capture())
}

func TestSelectFormField_TapFocus(t *testing.T) {
	test.NewApp()
	defer test.NewApp()
	test.ApplyTheme(t, test.Theme())
	sf := NewSelectFormField("Name", "Andrea", []string{"Peter", "Andrea"})
	sf.Placeholder = "Your name"
	sf.Hint = "A hint text"

	w := test.NewWindow(sf)
	w.Resize(fyne.NewSize(150, 180))
	defer w.Close()

	w.Canvas().Focus(sf.selectField)
	test.AssertImageMatches(t, "select_form_field/focused.png", w.Canvas().Capture())

	w.Canvas().Focus(nil)
	test.AssertImageMatches(t, "select_form_field/unfocused.png", w.Canvas().Capture())

	test.Tap(sf.selectField)
	test.AssertImageMatches(t, "select_form_field/focused_with_popup.png", w.Canvas().Capture())
	test.TapCanvas(w.Canvas(), fyne.NewPos(140, 170))
	test.AssertImageMatches(t, "select_form_field/focused.png", w.Canvas().Capture())

	r := test.WidgetRenderer(sf).(*formFieldRenderer)

	test.Tap(r.labelBg)
	test.AssertImageMatches(t, "select_form_field/focused_with_popup.png", w.Canvas().Capture())
	test.TapCanvas(w.Canvas(), fyne.NewPos(140, 170))
	test.AssertImageMatches(t, "select_form_field/focused.png", w.Canvas().Capture())
}

func TestSelectFormField_DisabledTryFocus(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	sf := NewSelectFormField("Name", "Andrea", []string{"Peter", "Andrea"})
	sf.Placeholder = "Your name"
	sf.Hint = "A hint text"
	sf.Disable()

	w := test.NewWindow(sf)
	w.Resize(fyne.NewSize(150, 180))
	defer w.Close()

	w.Canvas().Focus(sf.selectField)
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())
	w.Canvas().Focus(nil)

	r := test.WidgetRenderer(sf).(*formFieldRenderer)

	test.Tap(r.labelBg)
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())
	w.Canvas().Focus(nil)
}

func TestSelectFormField_DisabledTryHover(t *testing.T) {
	test.ApplyTheme(t, test.Theme())
	sf := NewSelectFormField("Name", "Andrea", []string{"Peter", "Andrea"})
	sf.Placeholder = "Your name"
	sf.Hint = "A hint text"
	sf.Disable()

	w := test.NewWindow(sf)
	w.Resize(fyne.NewSize(150, 180))
	defer w.Close()

	sf.selectField.MouseIn(&desktop.MouseEvent{})
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())

	sf.selectField.MouseOut()
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())

	r := test.WidgetRenderer(sf).(*formFieldRenderer)

	r.labelBg.MouseIn(&desktop.MouseEvent{})
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())

	r.labelBg.MouseOut()
	test.AssertImageMatches(t, "select_form_field/disabled.png", w.Canvas().Capture())
}
