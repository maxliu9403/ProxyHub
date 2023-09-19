package callAPi

import (
	"context"
	"fmt"
	"github.com/maxliu9403/common/gadget"
	"github.com/maxliu9403/common/httputil"
	"io/ioutil"
	"net/http"
)

const (
	HTTPGet  = "get"
	HTTPPost = "post"
)

func DoHTTP(ctx context.Context, action, url string, sendOptions ...httputil.SendOption) (respBytes []byte, err error) {
	var resp *http.Response

	spanCtx, err := gadget.ExtractTraceSpan(ctx)
	if err == nil {
		sendOptions = append(sendOptions, httputil.SendTraceCTX(spanCtx))
	}

	switch action {
	case HTTPGet:
		resp, err = httputil.Get(url, sendOptions...) //nolint:bodyclose
	case HTTPPost:
		resp, err = httputil.Post(url, sendOptions...) //nolint:bodyclose
	default:
		err = fmt.Errorf("unknown action %s", action)
	}

	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
