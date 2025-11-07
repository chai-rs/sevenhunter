package errx

import "errors"

type Error struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal error  `json:"-"`
}

func (e *Error) Error() string {
	return e.Internal.Error()
}

func E(code int, internal error, messages ...string) *Error {
	e := new(Error)
	e.Code = code
	e.Internal = internal

	if len(messages) > 0 {
		e.Message = messages[0]
	} else {
		e.Message = internal.Error()
	}

	return e
}

func M(code int, message string) *Error {
	return E(code, errors.New(message), message)
}
