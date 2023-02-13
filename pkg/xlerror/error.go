package xlerror

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	codes = map[int]struct{}{}
)

var (
	NullError       = add(0, "")
	Success         = add(200, "SUCCESS")
	ErrRequest      = add(400, "请求参数错误")
	ErrNotFind      = add(404, "没有找到")
	ErrForbidden    = add(403, "请求被拒绝")
	ErrNoPermission = add(405, "无权限")
	ErrServer       = add(500, "请稍后重试")
	ErrTimeout      = add(504, "请求超时")
	ErrRateLimit    = add(600, "请稍后重试")
	ErrForWait      = add(900, "请稍后重试")
	ErrToken        = add(1000, "错误的Token")
)

// New 创建一个错误
func New(code int, msg string) Error {
	if code < 1000 && code > 0 {
		panic("error code must be greater than 1000")
	}
	return add(code, msg)
}

// add only inner error
func add(code int, msg string) Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", code))
	}
	codes[code] = struct{}{}
	return Error{
		code: code, message: msg,
	}
}

type Errors interface {
	// Error sometimes Error return Code in string form
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
	// Details get error detail,it may be nil.
	Details() []interface{}
	// Equal for compatible.
	Equal(error) bool
	// Reload Message
	Reload(string) Error
}

type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Reload(message string) Error {
	e.message = message
	return e
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Details() []interface{} { return nil }

func (e Error) Equal(err error) bool { return Equal(err, e) }

func String(e string) Error {
	if e == "" {
		return NullError
	}
	return Error{
		code:    500,
		message: e,
	}
}

// Cause 解析错误码
func Cause(err error) Errors {
	if err == nil {
		return NullError
	}
	if ec, ok := errors.Cause(err).(Errors); ok {
		return ec
	}
	return ErrServer // String(err.Error())
}

// Equal 两个错误错误码是否一致
func Equal(err error, e Error) bool {
	return Cause(err).Code() == e.Code()
}

// Wrap 对错误描述进行包装
func Wrap(err error, message string) Errors {
	ec := Cause(err)
	return Error{
		code:    ec.Code(),
		message: fmt.Sprintf("%s: %s", ec.Message(), message),
	}
}
