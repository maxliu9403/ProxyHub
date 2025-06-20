// Copyright (c) 2016-2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package httputil

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/maxliu9403/common/gadget"
	"github.com/maxliu9403/common/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

const schemeHTTPS = `https`

// RoundTripper is an alias of the http.RoundTripper for mocking purposes.
type RoundTripper = http.RoundTripper

// StatusError occurs if an HTTP response has an unexpected status code.
type StatusError struct {
	Method       string
	URL          string
	Status       int
	Header       http.Header
	ResponseDump string
}

// NewStatusError returns a new StatusError.
func NewStatusError(resp *http.Response) StatusError {
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	respDump := string(respBytes)
	if err != nil {
		respDump = fmt.Sprintf("failed to dump response: %s", err)
	}

	return StatusError{
		Method:       resp.Request.Method,
		URL:          resp.Request.URL.String(),
		Status:       resp.StatusCode,
		Header:       resp.Header,
		ResponseDump: respDump,
	}
}

func (e StatusError) Error() string {
	if e.ResponseDump == "" {
		return fmt.Sprintf("%s %s %d", e.Method, e.URL, e.Status)
	}
	return fmt.Sprintf("%s %s %d: %s", e.Method, e.URL, e.Status, e.ResponseDump)
}

// IsStatus returns true if err is a StatusError of the given status.
func IsStatus(err error, status int) bool {
	statusErr, ok := err.(StatusError)
	return ok && statusErr.Status == status
}

// IsCreated returns true if err is a "created", 201
func IsCreated(err error) bool {
	return IsStatus(err, http.StatusCreated)
}

// IsNotFound returns true if err is a "not found" StatusError.
func IsNotFound(err error) bool {
	return IsStatus(err, http.StatusNotFound)
}

// IsConflict returns true if err is a "status conflict" StatusError.
func IsConflict(err error) bool {
	return IsStatus(err, http.StatusConflict)
}

// IsAccepted returns true if err is a "status accepted" StatusError.
func IsAccepted(err error) bool {
	return IsStatus(err, http.StatusAccepted)
}

// IsForbidden returns true if statis code is 403 "forbidden"
func IsForbidden(err error) bool {
	return IsStatus(err, http.StatusForbidden)
}

// NetworkError occurs on any Send error which occurred while trying to send
// the HTTP request, e.g. the given host is unresponsive.
type NetworkError struct {
	err error
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("network error: %s", e.err)
}

// IsNetworkError returns true if err is a NetworkError.
func IsNetworkError(err error) bool {
	_, ok := err.(NetworkError)
	return ok
}

type sendOptions struct {
	body          io.Reader
	timeout       time.Duration
	acceptedCodes map[int]bool
	headers       map[string]string
	redirect      func(req *http.Request, via []*http.Request) error
	retry         retryOptions
	transport     http.RoundTripper
	ctx           context.Context
	spanCtx       context.Context
	span          opentracing.Span
	// This is not a valid http option. It provides a way to override
	// parts of the url. For example, url.Scheme can be changed from
	// http to https.
	url *url.URL

	// This is not a valid http option. HTTP fallback is added to allow
	// easier migration from http to https.
	// In go1.11 and go1.12, the responses returned when http request is
	// sent to https server are different in the fallback mode:
	// go1.11 returns a network error whereas go1.12 returns BadRequest.
	// This causes TestTLSClientBadAuth to fail because the test checks
	// retry error.
	// This flag is added to allow disabling http fallback in unit tests.
	// NOTE: it does not impact how it runs in production.
	httpFallbackDisabled bool
}

// SendOption allows overriding defaults for the Send function.
type SendOption func(*sendOptions)

// SendNoop returns a no-op option.
func SendNoop() SendOption {
	// no-op
	return func(*sendOptions) {}
}

// SendBody specifies a body for http request
func SendBody(body io.Reader) SendOption {
	return func(o *sendOptions) { o.body = body }
}

// SendTimeout specifies timeout for http request
func SendTimeout(timeout time.Duration) SendOption {
	return func(o *sendOptions) { o.timeout = timeout }
}

// SendHeaders specifies headers for http request
func SendHeaders(headers map[string]string) SendOption {
	return func(o *sendOptions) {
		for k, v := range headers {
			o.headers[k] = v
		}
	}
}

// SendAcceptedCodes specifies accepted codes for http request
func SendAcceptedCodes(codes ...int) SendOption {
	m := make(map[int]bool)
	for _, c := range codes {
		m[c] = true
	}
	return func(o *sendOptions) { o.acceptedCodes = m }
}

// SendRedirect specifies a redirect policy for http request
func SendRedirect(redirect func(*http.Request, []*http.Request) error) SendOption {
	return func(o *sendOptions) { o.redirect = redirect }
}

type retryOptions struct {
	backoff    backoff.BackOff
	extraCodes map[int]bool
}

// RetryOption allows overriding defaults for the SendRetry option.
type RetryOption func(*retryOptions)

// RetryBackoff adds exponential backoff between retries.
func RetryBackoff(b backoff.BackOff) RetryOption {
	return func(o *retryOptions) { o.backoff = b }
}

// RetryCodes adds more status codes to be retried (in addition to the default
// 5XX codes).
//
// WARNING: You better know what you're doing to retry anything non-5XX.
func RetryCodes(codes ...int) RetryOption {
	return func(o *retryOptions) {
		for _, c := range codes {
			o.extraCodes[c] = true
		}
	}
}

// SendRetry will we retry the request on network / 5XX errors.
func SendRetry(options ...RetryOption) SendOption {
	retry := retryOptions{
		// backoff: backoff.WithMaxRetries(
		//	backoff.NewConstantBackOff(250*time.Millisecond),
		//	2),
		backoff:    backoff.NewConstantBackOff(250 * time.Millisecond),
		extraCodes: make(map[int]bool),
	}
	for _, o := range options {
		o(&retry)
	}
	return func(o *sendOptions) { o.retry = retry }
}

// DisableHTTPFallback disables http fallback when https request fails.
func DisableHTTPFallback() SendOption {
	return func(o *sendOptions) {
		o.httpFallbackDisabled = true
	}
}

// SendTLS sets the transport with TLS config for the HTTP client.
func SendTLS(config *tls.Config) SendOption {
	return func(o *sendOptions) {
		if config == nil {
			return
		}
		o.transport = &http.Transport{TLSClientConfig: config}
		o.url.Scheme = schemeHTTPS
	}
}

// SendTLSTransport sets the transport with TLS config for the HTTP client.
func SendTLSTransport(transport http.RoundTripper) SendOption {
	return func(o *sendOptions) {
		o.transport = transport
		o.url.Scheme = schemeHTTPS
	}
}

// SendTransport sets the transport for the HTTP client.
func SendTransport(transport http.RoundTripper) SendOption {
	return func(o *sendOptions) { o.transport = transport }
}

// SendContext sets the context for the HTTP client.
func SendContext(ctx context.Context) SendOption {
	return func(o *sendOptions) { o.ctx = ctx }
}

func SendTraceCTX(spanCtx context.Context) SendOption {
	return func(o *sendOptions) {
		o.spanCtx = spanCtx
	}
}

// Send sends an HTTP request. May return NetworkError or StatusError (see above).
func Send(method, rawurl string, options ...SendOption) (*http.Response, error) { //nolint:funlen,gocognit
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, fmt.Errorf("parse url failed: %s", err.Error())
	}
	opts := &sendOptions{
		body:                 nil,
		timeout:              60 * time.Second,
		acceptedCodes:        map[int]bool{http.StatusOK: true},
		headers:              map[string]string{},
		retry:                retryOptions{backoff: &backoff.StopBackOff{}},
		transport:            nil, // Use HTTP default.
		ctx:                  context.Background(),
		url:                  u,
		httpFallbackDisabled: false,
	}
	for _, o := range options {
		o(opts)
	}

	if opts.spanCtx != nil {
		opts.span, _ = opentracing.StartSpanFromContext(opts.spanCtx, fmt.Sprintf("%s_%s://%s%s", method, u.Scheme, u.Host, u.Path))
		defer opts.span.Finish()
	}

	req, err := newRequest(method, opts)
	if err != nil {
		return nil, err
	}

	if opts.span != nil {
		ext.SpanKindRPCClient.Set(opts.span)
		ext.HTTPUrl.Set(opts.span, rawurl)
		ext.HTTPMethod.Set(opts.span, method)
		_ = opts.span.Tracer().Inject(
			opts.span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
	}

	client := http.Client{
		Timeout:       opts.timeout,
		CheckRedirect: opts.redirect,
		Transport:     opts.transport,
	}

	var resp *http.Response
	for {
		resp, err = client.Do(req)
		// Retry without tls. During migration there would be a time when the
		// component receiving the tls request does not serve https response.
		// TODO (@evelynl): disable retry after tls migration.
		if err != nil && req.URL.Scheme == schemeHTTPS && !opts.httpFallbackDisabled {
			logger.Warnf("failed to send https request: %s. Retrying with http...", err)
			var httpReq *http.Request
			httpReq, err = newRequest(method, opts)
			if err != nil {
				if opts.span != nil {
					opts.span.LogFields(
						log.Int("code", resp.StatusCode),
						log.Error(err),
					)
				}
				return nil, err
			}
			httpReq.URL.Scheme = "http"
			resp, err = client.Do(httpReq)
		}
		if err != nil ||
			(resp.StatusCode >= 500 && !opts.acceptedCodes[resp.StatusCode]) ||
			(opts.retry.extraCodes[resp.StatusCode]) {
			d := opts.retry.backoff.NextBackOff()
			if d == backoff.Stop {
				break // Backoff timed out.
			}
			time.Sleep(d)
			continue
		}
		break
	}
	if err != nil {
		if opts.span != nil {
			opts.span.LogFields(
				log.Int("code", 000),
				log.Error(NetworkError{err}),
			)
		}
		return nil, NetworkError{err}
	}
	if !opts.acceptedCodes[resp.StatusCode] {
		if opts.span != nil {
			opts.span.LogFields(
				log.Int("code", resp.StatusCode),
				log.Error(nopCloserStatusError(resp)),
			)
		}
		return nil, NewStatusError(resp)
	}
	if opts.span != nil {
		body, _ := ioutil.ReadAll(resp.Body)

		fi2 := ioutil.NopCloser(bytes.NewBuffer(body))
		resp.Body = fi2

		opts.span.LogFields(
			log.Int("code", resp.StatusCode),
			log.String("response", string(body)),
		)
	}

	return resp, nil
}

// Get sends a GET http request.
func Get(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodGet, url, options...)
}

// Head sends a HEAD http request.
func Head(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodHead, url, options...)
}

// Post sends a POST http request.
func Post(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodPost, url, options...)
}

// Put sends a PUT http request.
func Put(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodPut, url, options...)
}

// Patch sends a PATCH http request.
func Patch(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodPatch, url, options...)
}

// Delete sends a DELETE http request.
func Delete(url string, options ...SendOption) (*http.Response, error) {
	return Send(http.MethodDelete, url, options...)
}

// PollAccepted wraps GET requests for endpoints which require 202-polling.
func PollAccepted(url string, b backoff.BackOff, options ...SendOption) (*http.Response, error) {
	b.Reset()
	for {
		resp, err := Get(url, options...)
		if err != nil {
			if IsAccepted(err) {
				d := b.NextBackOff()
				if d == backoff.Stop {
					break // Backoff timed out.
				}
				time.Sleep(d)
				continue
			}
			return nil, err
		}
		return resp, nil
	}
	return nil, errors.New("backoff timed out on 202 responses")
}

func SendWithCtx(ctx context.Context, method, url string, sendOptions ...SendOption) (respBytes []byte, err error) {
	spanCtx, err := gadget.ExtractTraceSpan(ctx)
	if err == nil {
		sendOptions = append(sendOptions, SendTraceCTX(spanCtx), SendContext(ctx))
	} else {
		sendOptions = append(sendOptions, SendContext(ctx))
	}

	resp, err := Send(method, url, sendOptions...)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func newRequest(method string, opts *sendOptions) (*http.Request, error) {
	req, err := http.NewRequest(method, opts.url.String(), opts.body)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err)
	}

	req = req.WithContext(opts.ctx)
	if opts.body == nil {
		req.ContentLength = 0
	}

	for key, val := range opts.headers {
		req.Header.Set(key, val)
	}

	return req, nil
}

func nopCloserStatusError(resp *http.Response) StatusError {
	body, err := ioutil.ReadAll(resp.Body)
	respDump := string(body)
	if err != nil {
		respDump = fmt.Sprintf("failed to dump response: %s", err)
	}
	fi2 := ioutil.NopCloser(bytes.NewBuffer(body))
	resp.Body = fi2

	return StatusError{
		Method:       resp.Request.Method,
		URL:          resp.Request.URL.String(),
		Status:       resp.StatusCode,
		Header:       resp.Header,
		ResponseDump: respDump,
	}
}
