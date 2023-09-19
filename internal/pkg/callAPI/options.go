package callAPi

import "time"

type callOptions struct {
	timeout time.Duration
	header  map[string]string
	method  string
}

type CallOption func(*callOptions)

func SetMethod(method string) CallOption {
	return func(o *callOptions) { o.method = method }
}

func SetHeader(header map[string]string) CallOption {
	return func(o *callOptions) { o.header = header }
}

func SetTimeout(timeOut time.Duration) CallOption {
	return func(o *callOptions) { o.timeout = timeOut }
}
