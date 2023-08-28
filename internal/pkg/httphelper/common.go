package httphelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/uuid"
)

// correlatedId gets Correlation ID from supplied context. If no Correlation ID header is
// present in the supplied context, one will be created along with a value.
func correlatedId(ctx context.Context) string {
	correlation := utils.FromContext(ctx, constants.CorrelationHeader)
	if len(correlation) == 0 {
		correlation = uuid.New().String()
	}
	return correlation
}

func langFromCtx(ctx context.Context) string {
	lang := utils.FromContext(ctx, constants.AcceptLanguage)
	if len(lang) == 0 {
		return ""
	}
	return lang
}

// Helper method to get the body from the response after making the request
func getBody(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, errort.NewCommonEdgeX(errort.DefaultSystemError, "failed to get the body from the response", err)
	}
	return body, nil
}

// Helper method to make the request and return the response
func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, errort.NewCommonEdgeX(errort.DefaultSystemError, "failed to send a http request", err)
	}
	return resp, nil
}

func createRequest(ctx context.Context, httpMethod string, baseUrl string, requestPath string, requestParams url.Values) (*http.Request, error) {
	u, err := url.Parse(baseUrl)
	if err != nil {
		return nil, errort.NewCommonEdgeX(errort.DefaultSystemError, "fail to parse baseUrl", err)
	}
	u.Path = requestPath
	if requestParams != nil {
		u.RawQuery = requestParams.Encode()
	}
	req, err := http.NewRequest(httpMethod, u.String(), nil)
	if err != nil {
		return nil, errort.NewCommonEdgeX(errort.DefaultSystemError, "failed to create a http request", err)
	}

	req.Header.Set(constants.CorrelationHeader, correlatedId(ctx))
	req.Header.Set(constants.AcceptLanguage, langFromCtx(ctx))
	return req, nil
}

func createRequestWithRawData(ctx context.Context, httpMethod string, url string, data interface{}) (*http.Request, error) {
	jsonEncodedData, err := json.Marshal(data)
	if err != nil {
		return nil, errort.NewCommonEdgeX(errort.DefaultSystemError, "failed to encode input data to JSON", err)
	}

	content := utils.FromContext(ctx, constants.ContentType)
	if content == "" {
		content = constants.ContentTypeJSON
	}

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(jsonEncodedData))
	if err != nil {
		return nil, errort.NewCommonEdgeX(errort.DefaultSystemError, "failed to create a http request", err)
	}
	req.Header.Set(constants.ContentType, content)
	req.Header.Set(constants.CorrelationHeader, correlatedId(ctx))
	req.Header.Set(constants.AcceptLanguage, langFromCtx(ctx))
	return req, nil
}

// sendRequest will make a request with raw data to the specified URL.
// It returns the body as a byte array if successful and an error otherwise.
func sendRequest(ctx context.Context, req *http.Request) ([]byte, error) {
	data := ctx.Value("data")
	if data != nil {
		switch data.(type) {
		case []byte:
			req.Header.Set("data", string(data.([]byte)))
		case string:
			req.Header.Set("data", data.(string))
		default:
			req.Header.Set("dpType", fmt.Sprintln(reflect.TypeOf(data)))
			marshal, _ := json.Marshal(data)
			req.Header.Set("data", string(marshal))
		}
	}
	resp, err := makeRequest(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errort.NewCommonEdgeX(errort.DefaultSystemError, "the response should not be a nil", nil)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed, code:%v", resp.StatusCode)
	}

	return getBody(resp)
}

func CommResToSpecial(in interface{}, out interface{}) error {
	b, err := json.Marshal(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &out)
	if err != nil {
		return err
	}
	return nil
}
