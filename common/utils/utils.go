package utils

import (
	"encoding/json"
	"net"
	"strconv"
)

func Marshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if nil != err {
		return "", err
	} else {
		return string(b), nil
	}
}

func StrToInt(str string) int {
	intVar, _ := strconv.Atoi(str)
	return intVar
}

func ParseIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func ParseCidr(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}
