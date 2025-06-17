package callAPi

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/maxliu9403/common/httputil"
	"github.com/maxliu9403/common/logger"
)

// CallAPI 默认发起POST请求
func CallAPI(ctx context.Context, url string, params, resp interface{}, options ...CallOption) (err error) {
	opts := &callOptions{}
	for _, o := range options {
		o(opts)
	}

	sendBody, err := json.Marshal(&params)
	if err != nil {
		return
	}

	header := map[string]string{"Content-Type": "application/json", "operator": "ProxyHub", "User-Agent": "ProxyHub"}
	if len(opts.header) != 0 {
		header = opts.header
	}

	timeOut := 3 * time.Second
	if opts.timeout != 0 {
		timeOut = opts.timeout
	}

	method := HTTPPost
	if len(opts.method) != 0 {
		method = opts.method
	}

	logger.Debugf("request url is %s, params: %s", url, string(sendBody))
	requestOpts := []httputil.SendOption{
		httputil.SendBody(bytes.NewReader(sendBody)),
		httputil.SendHeaders(header),
		httputil.SendTimeout(timeOut),
	}

	bodyByte, err := DoHTTP(ctx, method, url, requestOpts...)
	if err != nil {
		return
	}

	logger.Debugf("response of url %s is: %s", url, string(bodyByte))
	if err = json.Unmarshal(bodyByte, resp); err != nil {
		return
	}

	return err
}
