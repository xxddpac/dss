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

func GetIpListByCidr(subnet string) []string {
	var (
		ips = make([]string, 0)
	)
	ip, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil
	}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); include(ip) {
		ips = append(ips, ip.String())
	}
	return ips[1 : len(ips)-1]
}

func include(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
