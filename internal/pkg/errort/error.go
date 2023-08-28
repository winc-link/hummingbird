package errort

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

type EdgeX interface {
	// Error obtains the error message associated with the error.
	Error() string
	// DebugMessages returns a detailed string for debug purpose.
	DebugMessages() string
	// Message returns the first level error message without further details.
	Message() string
	// Code returns the status code of this error.
	Code() uint32
}

// CommonEdgeX generalizes an error structure which can be used for any type of EdgeX error.
type CommonEdgeX struct {
	// callerInfo contains information of function call stacks.
	callerInfo string
	// message contains detailed information about the error.
	message string
	// code is the status code to represent this error.
	// We are using the standard HTTP status code: https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html
	code uint32
	// err is a nested error which is used to form a chain of errors for better context.
	err error
}

func Is(code uint32, err error) bool {
	ce, ok := As(err)
	if !ok {
		return false
	}

	return ce.code == code
}

func As(err error) (*CommonEdgeX, bool) {
	var ce = new(CommonEdgeX)
	ok := errors.As(err, &ce)

	return ce, ok
}

func Kind(err error) uint32 {
	var e CommonEdgeX
	if !errors.As(err, &e) {
		return KindUnknown
	}
	// We want to return the first "Kind" we see that isn't Unknown, because
	// the higher in the stack the Kind was specified the more context we had.
	if e.code != KindUnknown || e.err == nil {
		return e.code
	}

	return e.code
}

// Error creates an error message taking all nested and wrapped errors into account.
func (ce CommonEdgeX) Error() string {
	if ce.err == nil {
		return "[Code(" + strconv.Itoa(int(ce.Code())) + ")] " + ce.message
	}

	// ce.err.Error functionality gets the error message of the nested error and which will handle both CommonEdgeX
	// types and Go standard errors(both wrapped and non-wrapped).
	if ce.message != "" {
		return ce.message + " -> " + ce.err.Error()
	}
	return ce.err.Error()
}

// DebugMessages returns a string taking all nested and wrapped operations and errors into account.
func (ce CommonEdgeX) DebugMessages() string {
	if ce.err == nil {
		return ce.callerInfo + ": " + ce.message
	}

	if w, ok := ce.err.(CommonEdgeX); ok {
		return ce.callerInfo + ": " + ce.message + " -> " + w.DebugMessages()
	} else {
		return ce.callerInfo + ": " + ce.message + " -> " + ce.err.Error()
	}
}

// Message returns the first level error message without further details.
func (ce CommonEdgeX) Message() string {
	if ce.message == "" && ce.err != nil {
		if w, ok := ce.err.(CommonEdgeX); ok {
			return w.Message()
		} else {
			return ce.err.Error()
		}
	}

	return ce.message
}

// Code returns the status code of this error.
func (ce CommonEdgeX) Code() uint32 {
	return ce.code
}

// Unwrap retrieves the next nested error in the wrapped error chain.
// This is used by the new wrapping and unwrapping features available in Go 1.13 and aids in traversing the error chain
// of wrapped errors.
func (ce CommonEdgeX) Unwrap() error {
	return ce.err
}

// Is determines if an error is of type CommonEdgeX.
// This is used by the new wrapping and unwrapping features available in Go 1.13 and aids the errors.Is function when
// determining is an error or any error in the wrapped chain contains an error of a particular type.
func (ce CommonEdgeX) Is(err error) bool {
	switch err.(type) {
	case CommonEdgeX:
		return true
	default:
		return false

	}
}

func (ce CommonEdgeX) Cause() error {
	return ce.err
}

// NewCommonErr 封装自定义错误
// 使用场景: 在底层调用第三方服务或包时需要封装自定义错误
func NewCommonErr(code uint32, wrappedError error) error {
	_, ok := As(wrappedError)
	if ok {
		// 已经是自定义错误则不在封装
		return wrappedError
	}

	return errors.WithStack(&CommonEdgeX{
		code: code,
		err:  wrappedError,
	})
}

// NewCommonEdgeX alias NewCommonErr, 建议直接使用 NewCommonErr
// Deprecated
func NewCommonEdgeX(code uint32, message string, wrappedError error) error {
	if wrappedError != nil {
		return NewCommonErr(code, fmt.Errorf("%s: %w", message, wrappedError))
	} else {
		return NewCommonErr(code, fmt.Errorf("%s", message))
	}
}

// NewCommonEdgeXWrapper 封装自定义错误,1.取出自定义错误,2.未知错误统一使用DefaultSystemError
// 使用场景: 最上层需要取自定义错误码时调用
// Deprecated
func NewCommonEdgeXWrapper(wrappedError error) CommonEdgeX {
	code := DefaultSystemError
	// 优先断言默认edge错误，否则直接转换rpc的错误
	if w, ok := As(wrappedError); ok {
		code = w.code
	} else {
		ew := ConvertFromRPC(wrappedError)
		if w, ok = As(ew); ok {
			code = w.code
		}
	}
	// 非内部定制错误,强制转换为系统错误
	if code < DefaultSystemError {
		code = DefaultSystemError
	}
	return CommonEdgeX{
		callerInfo: getCallerInformation(),
		message:    "",
		code:       code,
		err:        wrappedError,
	}
}

// getCallerInformation generates information about the caller function. This function skips the caller which has
// invoked this function, but rather introspects the calling function 3 frames below this frame in the call stack. This
// function is a helper function which eliminates the need for the 'callerInfo' field in the `CommonEdgeX` type and
// providing an 'callerInfo' string when creating an 'CommonEdgeX'
func getCallerInformation() string {
	pc := make([]uintptr, 10)
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	return fmt.Sprintf("[%s]-%s(line %d)", file, f.Name(), line)
}
