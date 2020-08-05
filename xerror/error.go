package xerror

import "fmt"

type Error struct {
	HttpStatus int    `json:"-"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Prompt     string `json:"prompt"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%v]:%v", e.Code, e.Msg)
}

func (e *Error) CopyWithPrompt(p string) *Error {
	return &Error{
		e.HttpStatus, e.Code, e.Msg, p,
	}
}

func (e *Error) As(err interface{}) bool {
	_, ok := err.(*Error)
	return ok
}

func (e *Error) Is(err interface{}) bool {
	rel, ok := err.(*Error)
	if !ok {
		return false
	}
	return e == rel
}

func register(status, code int, m, p string) *Error {
	return &Error{status, code, m, p}
}
