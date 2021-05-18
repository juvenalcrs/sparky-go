package svalid

import (
	"fyne.io/fyne/v2"
)

// NewGroup creates a new validator as a result of combined validators.
func NewGroup(validators ...fyne.StringValidator) fyne.StringValidator {
	return func(s string) error {
		for _, validator := range validators {
			if err := validator(s); err != nil {
				return err
			}
		}
		return nil
	}
}
