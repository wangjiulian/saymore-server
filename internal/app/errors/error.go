package errors

import "fmt"

type Error interface {
	Error() string
	Code() int
	Msg() string
}

// Wrap system error
func Wrap(err error) Error {
	return vsErr{
		code: CodeActionAbort,
		msg:  err.Error(),
	}
}

func WrapInvalid(msg string) Error {
	return vsErr{
		code: CodeInvalidParam,
		msg:  msg,
	}
}

func WrapParamTransfer(err error) Error {
	return vsErr{
		code: CodeParamTransfer,
		msg:  err.Error(),
	}
}

// WrapDetail with code and msg
func WrapDetail(code int, msg string) Error {
	return vsErr{
		code: code,
		msg:  msg,
	}
}

type vsErr struct {
	code int
	msg  string
}

func (e vsErr) Error() string {
	return fmt.Sprintf("err_code: %d, err_msg: %s", e.code, e.msg)
}

func (e vsErr) Code() int {
	return e.code
}

func (e vsErr) Msg() string {
	return e.msg
}
