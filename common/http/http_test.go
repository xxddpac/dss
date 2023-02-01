package http

import (
	"context"
	"fmt"
	"testing"
)

var (
	baseUrl = "https://httpbin.org/"
	request = NewClient(baseUrl)
)

func TestHttpGet(t *testing.T) {
	statusCode, byteStr, err := request.Get(context.Background(), "get")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(statusCode, string(byteStr))
}

func TestHttpPost(t *testing.T) {
	statusCode, byteStr, err := request.Post(context.Background(), "post", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(statusCode, string(byteStr))
}
