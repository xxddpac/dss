package wp

import (
	_ "embed"
	"strings"
)

//go:embed list
var __wpList string

// WeakPasswordList parser weak password from list
var WeakPasswordList = func() []string {
	var resp []string
	for _, item := range strings.Split(__wpList, "\n") {
		resp = append(resp, strings.TrimSpace(item))
	}
	return resp
}()
