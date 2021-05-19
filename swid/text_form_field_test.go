package swid

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/assert"
)

func TestTextFormField_EmptyValidator(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "")
	tf.Validator = func(s string) error {
		if s == "" {
			return errors.New("required")
		}
		return nil
	}

	w := test.NewWindow(tf)
	defer w.Close()

	test.AssertImageMatches(t, "text_form_field/empty_validator_initial.png", w.Canvas().Capture())

	w.Canvas().Focus(tf.textField)

	test.AssertImageMatches(t, "text_form_field/empty_validator_focused.png", w.Canvas().Capture())

	w.Canvas().Focus(nil)

	// TODO this needs to be updated when Fyne fixes this (see the underline color and no icon)
	test.AssertImageMatches(t, "text_form_field/empty_validator_unfocused.png", w.Canvas().Capture())
}

func TestTextFormField_Validation(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "wrong")
	tf.Validator = func(s string) error {
		if s == "wrong" {
			return errors.New("wrong")
		}
		return nil
	}
	tf.Hint = "A hint text"

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(100, 80))
	defer w.Close()

	t.Run("initial_error", func(t *testing.T) {
		// If the field is not empty and it has an error, then the error should be shown
		// from the beginning.
		test.AssertImageMatches(t, "text_form_field/invalid_unfocused.png", w.Canvas().Capture())
	})

	t.Run("focus_initial_error", func(t *testing.T) {
		w.Canvas().Focus(tf.textField)

		// The error hint should be hide, when the form field is focused.
		test.AssertImageMatches(t, "text_form_field/invalid_focused_cursor_start.png", w.Canvas().Capture())
	})

	t.Run("unfocus_initial_error", func(t *testing.T) {
		w.Canvas().Focus(nil)

		// Unfocusing the form field, makes the error hint appears again (only if the field is still invalid)
		test.AssertImageMatches(t, "text_form_field/invalid_unfocused.png", w.Canvas().Capture())
	})

	t.Run("valid_empty", func(t *testing.T) {
		// Validation no error but field empty (unfocused)
		tf.SetText("")
		test.AssertImageMatches(t, "text_form_field/valid_empty_unfocused.png", w.Canvas().Capture())
		// focus
		w.Canvas().Focus(tf.textField)
		test.AssertImageMatches(t, "text_form_field/valid_empty_focused.png", w.Canvas().Capture())
		// unfocus for next test
		w.Canvas().Focus(nil)
	})

	t.Run("valid_notempty", func(t *testing.T) {
		// Type valid text
		w.Canvas().Focus(tf.textField)
		test.Type(tf.textField, "wron")
		test.AssertImageMatches(t, "text_form_field/valid_notempty_focused.png", w.Canvas().Capture())

		// unfocus
		w.Canvas().Focus(nil)
		test.AssertImageMatches(t, "text_form_field/valid_notempty_unfocused.png", w.Canvas().Capture())
	})

	t.Run("invalid_no_initial", func(t *testing.T) {
		// Complete invalid text
		w.Canvas().Focus(tf.textField)
		test.Type(tf.textField, "g")
		test.AssertImageMatches(t, "text_form_field/invalid_focused_cursor_end.png", w.Canvas().Capture())

		// unfocus to see the error hint
		w.Canvas().Focus(nil)
		test.AssertImageMatches(t, "text_form_field/invalid_unfocused.png", w.Canvas().Capture())
	})
}

func TestTextFormField_Placeholder(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "")
	tf.Validator = func(s string) error {
		if s == "wrong" {
			return errors.New("wrong")
		}
		return nil
	}
	tf.Hint = "A hint text"
	tf.Placeholder = "Your name"

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	t.Run("placeholder_hidden_empty_unfocused", func(t *testing.T) {
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_hidden_empty_unfocused.png", w.Canvas().Capture())
	})

	t.Run("placeholder_visible_empty_focused", func(t *testing.T) {
		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "Your name", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_visible_empty_focused.png", w.Canvas().Capture())
	})

	t.Run("placeholder_hidden_notempty_unfocused", func(t *testing.T) {
		tf.SetText("Peter")
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_hidden_notempty_unfocused.png", w.Canvas().Capture())
	})

	t.Run("placeholder_hidden_notempty_focused", func(t *testing.T) {
		tf.SetText("Peter")
		w.Canvas().Focus(tf.textField)
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_hidden_notempty_focused.png", w.Canvas().Capture())
	})

	t.Run("placeholder_hidden_empty_unfocused_with_set", func(t *testing.T) {
		w.Canvas().Focus(nil)
		tf.SetText("")
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_hidden_empty_unfocused.png", w.Canvas().Capture())
	})

	t.Run("placeholder_visible_empty_focused_with_set", func(t *testing.T) {
		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		tf.SetText("")
		assert.Equal(t, "Your name", tf.Placeholder)
		assert.Equal(t, "Your name", tf.textField.PlaceHolder)
		test.AssertImageMatches(t,
			"text_form_field/placeholder_visible_empty_focused.png", w.Canvas().Capture())
	})
}

func TestTextFormField_SetText(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "")
	assert.Equal(t, "", tf.Text())
	assert.Equal(t, "", tf.textField.Text)

	tf.SetText("Peter")
	assert.Equal(t, "Peter", tf.Text())
	assert.Equal(t, "Peter", tf.textField.Text)
}

func TestTextFormField_Reset(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "Andrea")
	tf.Placeholder = "Your name"

	assert.Equal(t, "Andrea", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	tf.SetText("Peter")
	assert.Equal(t, "Peter", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	tf.Reset()
	assert.Equal(t, "Andrea", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	tf = NewTextFormField("Name", "")
	tf.Placeholder = "Your name"

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	assert.Equal(t, "", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	tf.SetText("Peter")
	assert.Equal(t, "Peter", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	tf.Reset()
	assert.Equal(t, "", tf.Text())
	assert.Equal(t, "", tf.textField.PlaceHolder)

	w.Canvas().Focus(tf.textField)
	assert.Equal(t, "", tf.Text())
	assert.Equal(t, "Your name", tf.textField.PlaceHolder)
}

func TestTextFormField_DisableValid(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "")
	tf.Placeholder = "Your name"

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	t.Run("disable_valid_empty_initial", func(t *testing.T) {
		tf.SetText("")
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_empty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_empty.png", w.Canvas().Capture())
	})

	t.Run("disable_valid_notempty", func(t *testing.T) {
		tf.SetText("Andrea")
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_notempty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_notempty.png", w.Canvas().Capture())
	})

	t.Run("disable_valid_empty", func(t *testing.T) {
		tf.SetText("")
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_empty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_valid_empty.png", w.Canvas().Capture())
	})
}

func TestTextFormField_DisableInvalid(t *testing.T) {
	test.ApplyTheme(t, theme.LightTheme())
	tf := NewTextFormField("Name", "")
	tf.Validator = func(s string) error {
		if s == "" {
			return errors.New("required")
		}
		if s == "wrong" {
			return errors.New("wrong")
		}
		return nil
	}
	tf.Placeholder = "Your name"

	w := test.NewWindow(tf)
	w.Resize(fyne.NewSize(130, 80))
	defer w.Close()

	t.Run("disable_invalid_empty_initial", func(t *testing.T) {
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_empty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_empty.png", w.Canvas().Capture())
	})

	t.Run("disable_invalid_empty", func(t *testing.T) {
		tf.SetText("")
		assert.True(t, tf.dirty)
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_empty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_empty.png", w.Canvas().Capture())
	})

	t.Run("disable_invalid_notempty", func(t *testing.T) {
		tf.SetText("wrong")
		tf.Disable()
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_notempty.png", w.Canvas().Capture())

		w.Canvas().Focus(tf.textField)
		defer w.Canvas().Focus(nil)
		assert.True(t, tf.Disabled())
		test.AssertImageMatches(t, "text_form_field/disable_invalid_notempty.png", w.Canvas().Capture())
	})
}
