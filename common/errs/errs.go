package errs

import (
	"errors"
	"fmt"
	"temp/common/errs/errorcode"
)

type RpcError struct {
	Code    uint32
	Message string
	Data    interface{}
}

func (e *RpcError) Error() string {
	return fmt.Sprintf("code: %d, message: %s, data: %v", e.Code, e.Message, e.Data)
}

// NewRpcError 创造新错误
func NewRpcError(code uint32, message string, data interface{}) *RpcError {
	return &RpcError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// New 创建已知错误
func New(e errorcode.ErrorCode, data interface{}) *RpcError {
	return NewRpcError(e.Code, e.Msg, data)
}

// Wrap 创建错误，支持 Wrap(err) 和 Wrap(code, err, detail) 两种形式。
func Wrap(args ...interface{}) *RpcError {
	if len(args) == 0 {
		return NewRpcError(errorcode.ServerError.Code, errorcode.ServerError.Msg, nil)
	}

	if len(args) == 1 {
		if err, ok := args[0].(error); ok {
			if err == nil {
				return nil
			}
			return NewRpcError(errorcode.ServerError.Code, errorcode.ServerError.Msg, err.Error())
		}
		if code, ok := args[0].(errorcode.ErrorCode); ok {
			return NewRpcError(code.Code, code.Msg, nil)
		}
		return NewRpcError(errorcode.ServerError.Code, errorcode.ServerError.Msg, fmt.Sprint(args[0]))
	}

	if code, ok := args[0].(errorcode.ErrorCode); ok {
		if err, ok := args[1].(error); ok {
			if err == nil {
				return NewRpcError(code.Code, code.Msg, nil)
			}
			if len(args) > 2 && args[2] != nil {
				return NewRpcError(code.Code, code.Msg, args[2])
			}
			return NewRpcError(code.Code, code.Msg, err.Error())
		}
		if len(args) > 1 && args[1] != nil {
			return NewRpcError(code.Code, code.Msg, args[1])
		}
		return NewRpcError(code.Code, code.Msg, nil)
	}

	if err, ok := args[0].(error); ok {
		if err == nil {
			return nil
		}
		return NewRpcError(errorcode.ServerError.Code, errorcode.ServerError.Msg, err.Error())
	}

	return NewRpcError(errorcode.ServerError.Code, errorcode.ServerError.Msg, fmt.Sprint(args...))
}

// From 获取错误
func From(err error) (*RpcError, bool) {
	if err == nil {
		return nil, false
	}
	var rpcErr *RpcError
	if errors.As(err, &rpcErr) {
		return rpcErr, true
	}
	return nil, false
}

// Message 格式化错误信息
func Message(e *RpcError) string {
	if e == nil {
		return ""
	}

	codeStr := fmt.Sprintf("%d", e.Code)
	dataStr := fmt.Sprintf("%#v", e.Data)

	return fmt.Sprintf("%s|%s|%s", codeStr, e.Message, dataStr)

}
