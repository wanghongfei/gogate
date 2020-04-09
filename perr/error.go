package perr

import (
	"errors"
	"fmt"
	"runtime"
)

// 增强内置error接口
type EnvError interface {
	error

	// 携带环境信息的错误描述
	ErrorWithEnv() string
}

// 不能暴露给用户的系统级错误
type SystemError struct {
	SrcName string
	LineNumber int

	Msg string

	Cause error
}

func (s SystemError) Error() string {
	return s.Msg
}

func (s SystemError) ErrorWithEnv() string {
	return fmt.Sprintf("[%s:%d] %s", s.SrcName, s.LineNumber, s.Msg)
}

func (s SystemError) Unwrap() error {
	return s.Cause
}


// 可以给用户看的业务错误
type BizError struct {
	SrcName string
	LineNumber int

	Msg string

	Cause error
}

func (b BizError) Error() string {
	return b.Msg
}

func (b BizError) ErrorWithEnv() string {
	return fmt.Sprintf("[%s:%d]%s", b.SrcName, b.LineNumber, b.Msg)
}

func (b BizError) Unwrap() error {
	return b.Cause
}

// 创建一个携带环境信息的error, 包含文件名 + 行号
func BizErrorf(format string, args ...interface{}) EnvError {
	_, srcName, line, _ := runtime.Caller(1)
	fmtErr := fmt.Errorf(format, args...)

	return &BizError{
		SrcName:    srcName,
		LineNumber: line,
		Msg:        fmtErr.Error(),
		Cause:		errors.Unwrap(fmtErr),
	}
}

// 创建一个携带环境信息的error, 包含文件名 + 行号
func SystemErrorf(format string, args ...interface{}) error {
	_, srcName, line, _ := runtime.Caller(1)
	fmtErr := fmt.Errorf(format, args...)


	return &SystemError{
		SrcName:    srcName,
		LineNumber: line,
		Msg:        fmtErr.Error(),
		Cause:		errors.Unwrap(fmtErr),
	}
}

func ParseError(err error) (*BizError, *SystemError, error) {
	// 是否为系统错误
	var sysErr *SystemError
	if errors.As(err, &sysErr) {
		return nil, sysErr, nil
	}

	// 是业务错误
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr, nil, nil
	}

	return nil, nil, err
}
