package wp

import (
	_ "embed"
	"strings"
)

//go:embed list
var __wpList string
var usernameList = []string{"root", "mysql", "redis"}

type UserPass struct {
	UserName string
	Password string
}

// WeakUserPassList parser weak password from list
var WeakUserPassList = func() []UserPass {
	var resp []UserPass
	for _, username := range usernameList {
		for _, password := range strings.Split(__wpList, "\n") {
			resp = append(resp, UserPass{
				UserName: username,
				Password: strings.TrimSpace(password),
			})
		}
	}
	return resp
}()
