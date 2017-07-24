// HTTP Client which handles generic error routing and marshaling
package httpclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/ContainX/depcon/pkg/encoding"
	"github.com/ContainX/depcon/pkg/logger"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

var log = logger.GetLogger("client")

type Response struct {
	Status  int
	Content string
	Elapsed time.Duration
	Error   error
}

type Request struct {
	// Http Method type
	method Method
	// Complete URL including params
	url string
	// Post data
	data string
	// Expected data type
	result interface{}
	// encoding type (optional : default JSON)
	encodingType encoding.EncoderType
}

type HttpClientConfig struct {
	sync.RWMutex
	// Http Basic Auth Username
	HttpUser string
	// Http Basic Auth Password
	HttpPass string
	// Http Authorization Token
	HttpToken string
	// Request timeout
	RequestTimeout int
	// TLS Insecure Skip Verify
	TLSInsecureSkipVerify bool
}

type HttpClient struct {
	config HttpClientConfig
	http   *http.Client
}

var (
	// invalid or error response
	ErrorInvalidResponse = errors.New("Invalid response from Remote")
	// some resource does not exists
	ErrorNotFound = errors.New("The resource does not exist")
	// Generic Error Message
	ErrorMessage = errors.New("Unknown error message was captured")
	// Not Authorized
	ErrorNotAuthorized = errors.New("Not Authorized to perform this action - Status: 403")
	// Not Authenticated
	ErrorNotAuthenticated = errors.New("Not Authenticated to perform this action - Status: 401")
)

func NewDefaultConfig() *HttpClientConfig {
	return &HttpClientConfig{RequestTimeout: 30, TLSInsecureSkipVerify: false}
}

func DefaultHttpClient() *HttpClient {
	return NewHttpClient(*NewDefaultConfig())
}

func NewHttpClient(config HttpClientConfig) *HttpClient {
	hc := &HttpClient{
		config: config,
		http: &http.Client{
			Timeout: time.Duration(config.RequestTimeout) * time.Second,
		},
	}
	if config.TLSInsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		hc.http.Transport = tr
	}
	return hc
}

func NewResponse(status int, elapsed time.Duration, content string, err error) *Response {
	return &Response{Status: status, Elapsed: elapsed, Content: content, Error: err}
}

func (h *HttpClient) HttpGet(url string, result interface{}) *Response {
	return h.invoke(&Request{method: GET, url: url, result: result})
}

func (h *HttpClient) HttpPut(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(PUT, url, data, result)
}

func (h *HttpClient) HttpDelete(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(DELETE, url, data, result)
}

func (h *HttpClient) HttpPost(url string, data interface{}, result interface{}) *Response {
	return h.httpCall(POST, url, data, result)
}

func (h *HttpClient) httpCall(method Method, url string, data interface{}, result interface{}) *Response {
	var body string
	if data != nil {
		body = h.convertBody(data)
	}

	r := &Request{
		method: method,
		url:    url,
		data:   body,
		result: result,
	}

	return h.invoke(r)
}

// Creates a net/http Request and associates default headers and authentication
// parameters
func (h *HttpClient) CreateHttpRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	AddDefaultHeaders(request)
	AddAuthentication(h.config, request)

	return request, nil
}

func (h *HttpClient) invoke(r *Request) *Response {

	log.Debug("%s - %s, Body:\n%s", r.method.String(), r.url, r.data)

	request, err := h.CreateHttpRequest(r.method.String(), r.url, strings.NewReader(r.data))

	if err != nil {
		return &Response{Error: err}
	}

	req_start := time.Now()
	response, err := h.http.Do(request)
	req_elapsed := time.Now().Sub(req_start)

	if err != nil {
		return NewResponse(0, req_elapsed, "", err)
	}

	status := response.StatusCode
	var content string
	if response.ContentLength != 0 {
		defer response.Body.Close()
		rc, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return NewResponse(status, req_elapsed, "", err)
		}
		content = string(rc)
	}

	log.Debug("Status: %v, RAW: %s", status, content)

	if status >= 200 && status < 300 {
		if r.result != nil {
			h.convert(r, content)
		}
		return NewResponse(status, req_elapsed, content, nil)
	}

	switch status {
	case 500:
		return NewResponse(status, req_elapsed, content, ErrorInvalidResponse)
	case 404:
		return NewResponse(status, req_elapsed, content, ErrorNotFound)
	case 403:
		return NewResponse(status, req_elapsed, content, ErrorNotAuthorized)
	case 401:
		return NewResponse(status, req_elapsed, content, ErrorNotAuthenticated)
	}

	return NewResponse(status, req_elapsed, content, ErrorMessage)
}

func (h *HttpClient) convertBody(data interface{}) string {
	if data == nil {
		return ""
	}
	encoder, _ := encoding.NewEncoder(encoding.JSON)
	body, _ := encoder.Marshal(data)
	return body
}

func (h *HttpClient) convert(r *Request, content string) error {
	um, _ := encoding.NewEncoder(encoding.JSON)
	if r.encodingType != 0 {
		um, _ = encoding.NewEncoder(r.encodingType)
	}
	um.UnMarshalStr(content, r.result)
	return nil

}

func AddDefaultHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
}

func AddAuthentication(c HttpClientConfig, req *http.Request) {
	if c.HttpToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token=%v", c.HttpToken))
		return
	}
	if c.HttpUser != "" {
		req.SetBasicAuth(c.HttpUser, c.HttpPass)
	}
}

func (h *HttpClient) Unwrap() *http.Client {
	return h.http
}

func (h *HttpClient) Configuration() HttpClientConfig {
	return h.config
}
