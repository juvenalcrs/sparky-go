package svalid

var errMsgs = ErrorMessages{
	NotEmpty:  "This field cannot be empty",
	Email:     "",
	MinLength: "Min length must be %d",
}

// ErrorMessages defines all the error messages.
type ErrorMessages struct {
	NotEmpty  string
	Email     string
	MinLength string
}

// ConfigErrMessages configure the error messages for validation.
func ConfigErrMessages(msgs *ErrorMessages) {
	if msgs.NotEmpty != "" {
		errMsgs.NotEmpty = msgs.NotEmpty
	}
	if msgs.Email != "" {
		errMsgs.Email = msgs.Email
	}
	if msgs.MinLength != "" {
		errMsgs.MinLength = msgs.MinLength
	}
}
