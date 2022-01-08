package boot

type WebError interface {
	Msg() string
	Code() string
	Error() string
}

func NewWebError(code string, msg string) WebError {
	return &webError{code, msg}
}

type webError struct {
	code string
	msg  string
}

func (e *webError) Msg() string {
	return e.msg
}

func (e *webError) Code() string {
	return e.code
}

func (e *webError) Error() string {
	return e.code + ":" + e.Msg()
}
