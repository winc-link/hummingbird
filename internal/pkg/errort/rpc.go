package errort

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRPCStatusErr 转换为GRPC错误
// TODO 已经添加 interceptor 统一处理不需要单独调用
func NewRPCStatusErr(err error) error {
	if err == nil {
		return nil
	}

	errw, ok := As(err)
	if !ok {
		return status.New(codes.Code(DefaultSystemError), err.Error()).Err()
	}
	st := status.New(codes.Code(errw.Code()), errw.Error())
	return st.Err()
}

// ConvertFromRPC 转换GRPC错误码
// TODO 已经添加 interceptor 统一处理不需要单独调用
func ConvertFromRPC(err error) error {
	if err == nil {
		return nil
	}

	st := status.Convert(err)
	if st == nil {
		return err
	}

	if st.Code() == codes.Unknown {
		return NewCommonErr(DefaultSystemError, err)
	}
	return NewCommonErr(uint32(st.Code()), err)
}
