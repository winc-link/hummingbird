package errort

import (
	"errors"
	"testing"
)

func TestErr(t *testing.T) {
	err := errors.New("this is err")
	errWarp := NewCommonEdgeXWrapper(err)
	errWarp = NewCommonEdgeXWrapper(errWarp)
	errWarp = NewCommonEdgeXWrapper(errWarp)
	errWarp = NewCommonEdgeXWrapper(errWarp)

	t.Logf("printferror: %v", errWarp)
	t.Log("errors.Error():", errWarp.Error())
	t.Log("code:", errWarp.Code())
	t.Log("message:", errWarp.Message())
	t.Log("debugMessages err:", errWarp.DebugMessages())
}
