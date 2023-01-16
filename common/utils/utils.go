package utils

import "encoding/json"

func Marshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if nil != err {
		return "", err
	} else {
		return string(b), nil
	}
}
