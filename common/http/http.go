package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Header map[string]string

func defaultHeader() Header {
	return Header{
		"Content-Type": "application/json;charset=utf-8",
	}
}

func NewClient(baseUrl string, headers ...Header) *Client {
	header := defaultHeader()
	if 0 != len(headers) {
		for key, val := range headers[0] {
			header[key] = val
		}
	}
	return &Client{
		baseUrl,
		header,
		newHttpClient(),
		&sync.Mutex{},
	}
}

type Client struct {
	baseUrl    string
	header     Header
	httpclient *http.Client
	lock       *sync.Mutex
}

func newHttpClient() *http.Client {
	transport := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		ExpectContinueTimeout: 2 * time.Second,
		ResponseHeaderTimeout: 120 * time.Second,
		DisableKeepAlives:     false,
		DisableCompression:    false,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConnsPerHost:   200,
		MaxIdleConns:          2000,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   100 * time.Second,
	}
}

func (c *Client) Get(ctx context.Context, uri string, headers ...Header) (int, []byte, error) {
	return c.Request(ctx, http.MethodGet, uri, nil, headers...)
}

func (c *Client) Post(ctx context.Context, uri string, body interface{}, headers ...Header) (int, []byte, error) {
	return c.Request(ctx, http.MethodPost, uri, body, headers...)
}

func (c *Client) Request(ctx context.Context, method, uri string, body interface{}, headers ...Header) (int, []byte, error) {
	req, err := c.newHttpRequest(ctx, method, uri, body, headers...)
	if err != nil {
		return -1, nil, fmt.Errorf("failed build http request, err:%s", err.Error())
	}
	resp, err := c.httpclient.Do(req)
	if err != nil {
		return -1, nil, fmt.Errorf("failed send http request, err: %s", err.Error())
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			panic(err)
		}
	}(resp.Body)
	rBytes, err := c.bytes(resp)
	return resp.StatusCode, rBytes, err
}

func (c *Client) newHttpRequest(ctx context.Context, method, uri string, body interface{}, headers ...Header) (*http.Request, error) {
	defer c.lock.Unlock()
	c.lock.Lock()
	var err error
	req := &http.Request{
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
	req = req.WithContext(ctx)
	req.URL, err = c.getUrl(uri)
	if err != nil {
		return nil, fmt.Errorf("failed get request uri, err: %s", err.Error())
	}
	for k, v := range c.getHeaders(headers...) {
		req.Header.Set(k, v)
	}
	if body != nil {
		if err = c.setBody(req, body); err != nil {
			return nil, fmt.Errorf("failed set request body, err: %s", err.Error())
		}
	}
	return req, nil
}

func (c *Client) setBody(req *http.Request, body interface{}) error {
	var (
		err    error
		buffer []byte
	)
	if val, ok := body.(io.Reader); ok {
		req.Body = io.NopCloser(val)
		return nil
	}
	if val, ok := body.(string); ok {
		buffer = []byte(val)
	} else if body != nil {
		buffer, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}
	req.Body = io.NopCloser(bytes.NewReader(buffer))
	req.ContentLength = int64(len(buffer))
	return nil
}

func (c *Client) getHeaders(headers ...Header) Header {
	if 0 != len(headers) {
		for key, val := range headers[0] {
			c.header[key] = val
		}
	}
	return c.header
}

func (c *Client) getUrl(uri string) (*url.URL, error) {
	if c.baseUrl == "" { // base url is empty, use uri as url
		return url.Parse(uri)
	}
	if "" == uri || strings.HasPrefix(uri, "/") {
		return url.Parse(c.baseUrl + uri)
	}
	return url.Parse(c.baseUrl + "/" + uri)
}

func (c *Client) bytes(resp *http.Response) ([]byte, error) {
	if strings.EqualFold("gzip", resp.Header.Get("Content-Encoding")) {
		reader, err := gzip.NewReader(resp.Body)
		if nil != err {
			return nil, err
		}
		return io.ReadAll(reader)
	}
	return io.ReadAll(resp.Body)
}
