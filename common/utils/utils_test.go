package utils

import (
	"fmt"
	"testing"
)

func TestParseIP(t *testing.T) {
	var (
		ip = "1.1.1.1"
	)
	fmt.Println(ParseIP(ip))
}

func TestParseCidr(t *testing.T) {
	var (
		cidr = "192.168.1.0/20"
	)
	fmt.Println(ParseCidr(cidr))
}
