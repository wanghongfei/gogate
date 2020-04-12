package perr

import (
	"errors"
	"fmt"
	"runtime"
)

// 增强内置error接口
type EnvError interface {
	error

	LineNumber() int
	SrcName() string
}

type WithBottomMsg interface {
	BottomMsg() string
}

// 不能暴露给用户的系统级错误
type SystemError struct {
	srcName string
	lineNumber int

	Msg string

	Cause error
}

func (s *SystemError) LineNumber() int {
	return s.lineNumber
}

func (s *SystemError) SrcName() string {
	return s.srcName
}

func (s *SystemError) Error() string {
	return s.Msg
}

func (s *SystemError) Unwrap() error {
	return s.Cause
}

// 可以给用户看的业务错误
type BizError struct {
	srcName string
	lineNumber int

	Msg string
	bottomMsg string

	Cause error
}

func (b *BizError) Error() string {
	return b.Msg
}

func (b *BizError) LineNumber() int {
	return b.lineNumber
}

func (b *BizError) SrcName() string {
	return b.srcName
}

func (b *BizError) BottomMsg() string {
	return b.bottomMsg
}

func (b *BizError) Unwrap() error {
	return b.Cause
}

// 创建一个携带环境信息的error, 包含文件名 + 行号
func WrapBizErrorf(cause error, format string, args ...interface{}) EnvError {
	_, srcName, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf(format, args...)

	bottomMsg := msg
	if nil != cause {
		if withBottomMsg, ok := cause.(WithBottomMsg); ok {
			bottomMsg = withBottomMsg.BottomMsg()
		}
	}

	return &BizError{
		srcName:    srcName,
		lineNumber: line,
		Msg:        msg,
		Cause:		cause,
		bottomMsg:  bottomMsg,
	}
}

// 创建一个携带环境信息的error, 包含文件名 + 行号
func WrapSystemErrorf(cause error, format string, args ...interface{}) error {
	_, srcName, line, _ := runtime.Caller(1)
	fmtErr := fmt.Errorf(format, args...)

	return &SystemError{
		srcName:    srcName,
		lineNumber: line,
		Msg:        fmtErr.Error(),
		Cause:		cause,
	}
}

func EnvMsg(err error) string {
	msg := err.Error()
	if envError, ok := err.(EnvError); ok {
		msg = fmt.Sprintf("[%s:%d]%s", envError.SrcName(), envError.LineNumber(), msg)
	}

	cause := err
	for {
		cause = errors.Unwrap(cause)
		if nil == cause {
			break
		}

		if envError, ok := cause.(EnvError); ok {
			msg = fmt.Sprintf("%s => [%s:%d]%s", msg, envError.SrcName(), envError.LineNumber(), cause.Error())
		}

	}

	return msg
}


func ParseError(err error) (*BizError, *SystemError, error) {
	// 是否包含系统错误
	var sysErr *SystemError
	if errors.As(err, &sysErr) {
		return nil, sysErr, nil
	}

	// 是否包含业务错误
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return bizErr, nil, nil
	}

	return nil, nil, err
}
