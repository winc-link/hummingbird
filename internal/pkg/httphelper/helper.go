package httphelper

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/i18n"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RenderFail(ctx context.Context, err error, w http.ResponseWriter, lc logger.LoggingClient) {
	lc.Errorf("renderFail: %v", err)
	errw := errort.NewCommonEdgeXWrapper(err)

	resp := CommonResponse{
		Result:    []interface{}{},
		Success:   false,
		ErrorMsg:  i18n.TransCode(ctx, errw.Code(), nil),
		ErrorCode: errw.Code(),
	}

	encode(resp, w)
}

func RenderFailNoLog(ctx context.Context, err error, w http.ResponseWriter) {
	errw := errort.NewCommonEdgeXWrapper(err)

	resp := CommonResponse{
		Result:    []interface{}{},
		Success:   false,
		ErrorMsg:  i18n.TransCode(ctx, errw.Code(), nil),
		ErrorCode: errw.Code(),
	}
	encode(resp, w)
}

type ResPageResult struct {
	List     interface{} `json:"list"`
	Total    uint32      `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

func Result(code uint32, data interface{}, msg string, w http.ResponseWriter, lc logger.LoggingClient) {
	success := false
	if code == errort.DefaultSuccess {
		success = true
	}
	if data == nil {
		data = []interface{}{}
	}
	resp := CommonResponse{
		Result:    data,
		Success:   success,
		ErrorMsg:  msg,
		ErrorCode: code,
	}
	encode(resp, w)
}

func ResultSuccess(data interface{}, w http.ResponseWriter, lc logger.LoggingClient) {
	Result(errort.DefaultSuccess, data, "success", w, lc)
}

func ResultNoLog(code uint32, data interface{}, msg string, w http.ResponseWriter) {
	success := false
	if code == errort.DefaultSuccess {
		success = true
	}
	if data == nil {
		data = []interface{}{}
	}
	resp := CommonResponse{
		Result:    data,
		Success:   success,
		ErrorMsg:  msg,
		ErrorCode: code,
	}
	encode(resp, w)
}

func ResultSuccessNoLog(data interface{}, w http.ResponseWriter) {
	ResultNoLog(errort.DefaultSuccess, data, "success", w)
}

func ResultZipFile(c *gin.Context, lc logger.LoggingClient, fileName string, file *os.File) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)

	defer file.Close()
	defer zw.Close()

	writer, err := zw.Create(fileName)
	if err != nil {
		lc.Errorf("failed to zip create header %v", err)
		c.JSON(http.StatusOK, NewFailCommonResponse(err))
		return
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		lc.Errorf("failed to  io copy %v", err)
		c.JSON(http.StatusOK, NewFailCommonResponse(err))
		return
	}

	// 主动关闭，将buf的zip内容写完
	zw.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", fileName))
	c.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func ResultExcelData(c *gin.Context, fileName string, data *bytes.Buffer) {
	c.Header("Response-Type", "blob")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", data.Bytes())
}

func NewPageResult(responses interface{}, total uint32, page int, pageSize int) ResPageResult {
	if responses == nil {
		responses = make([]interface{}, 0)
	}
	return ResPageResult{
		List:     responses,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}

type CommonResponse struct {
	Success    bool        `json:"success"`                     // 接口是否成功
	ErrorCode  uint32      `json:"errorCode"`                   // 错误码
	ErrorMsg   string      `json:"errorMsg,omitempty"`          // 错误信息
	SuccessMsg string      `json:"successMsg,omitempty"`        // 成功信息
	Result     interface{} `json:"result" swaggertype:"object"` // 返回结果
}

func (r CommonResponse) Error() error {
	if r.Success {
		return nil
	}
	return errort.NewCommonErr(r.ErrorCode, fmt.Errorf(r.ErrorMsg))
}

// deprecated
func NewFailCommonResponse(err error) CommonResponse {
	errw := errort.NewCommonEdgeXWrapper(err)

	return CommonResponse{
		Success:   false,
		ErrorMsg:  err.Error(),
		ErrorCode: errw.Code(),
		Result:    []interface{}{},
	}
}

func NewFailWithI18nResponse(ctx context.Context, err error) CommonResponse {
	errw := errort.NewCommonEdgeXWrapper(err)

	return CommonResponse{
		Success:   false,
		ErrorMsg:  i18n.TransCode(ctx, errw.Code(), nil),
		ErrorCode: errw.Code(),
		Result:    []interface{}{},
	}
}

func NewSuccessCommonResponse(data interface{}) CommonResponse {
	if data == nil {
		data = []interface{}{}
	}
	return CommonResponse{
		Success:   true,
		ErrorMsg:  "success",
		ErrorCode: errort.DefaultSuccess,
		Result:    data,
	}
}

// deprecated
func encode(i interface{}, w http.ResponseWriter) {
	w.Header().Add(constants.ContentType, constants.ContentTypeJSON)

	enc := json.NewEncoder(w)
	err := enc.Encode(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////
// websocket response Render
////////////////////////////////////////////////////////////////////////////////////////////////

func WsResultFail(code uint32, msg string) CommonResponse {
	resp := CommonResponse{
		Result:    []interface{}{},
		Success:   false,
		ErrorMsg:  msg,
		ErrorCode: code,
	}
	return resp
}

func WsResult(code uint32, data interface{}, errMsg string, successMsg string) CommonResponse {
	if data == nil {
		data = []interface{}{}
	}
	isSuccess := true
	if code != errort.DefaultSuccess {
		isSuccess = false
	}
	resp := CommonResponse{
		Result:     data,
		Success:    isSuccess,
		ErrorCode:  code,
		ErrorMsg:   errMsg,
		SuccessMsg: successMsg,
	}
	return resp
}

////////////////////////////////////////////////////////////////////////////////////////////////
// rpc response Render
////////////////////////////////////////////////////////////////////////////////////////////////

func RenderRpcFail(err error, lc logger.LoggingClient) error {
	lc.Errorf("renderRpcFail: %+v", err)

	errw := errort.NewCommonEdgeXWrapper(err)
	st := status.New(codes.Code(errw.Code()), err.Error())
	return st.Err()
}
