package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeLayout = "2006-01-02 15:04:05"
	TimeLayout        = "2006-01-02"
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

func UnixToString(timestamp int64) string {
	return time.Unix(timestamp, 0).Local().Format(defaultTimeLayout)
}

func StrToBool(str string) bool {
	b, _ := strconv.ParseBool(str)
	return b
}

func IsStrExists(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func WorkingDirectory() string {
	workingDirectory, _ := os.Getwd()
	return workingDirectory
}

func ParsePortRange(val string) (int, int, bool) {
	var (
		portStart, portEnd int
	)
	portRange := strings.Split(val, "-")
	if len(portRange) == 2 {
		portStart = StrToInt(portRange[0])
		portEnd = StrToInt(portRange[1])
	}
	return portStart, portEnd, len(portRange) == 2
}

func ParseIpRange(val string) (string, int, int, bool) {
	var (
		startIp                            string
		startIpEndSuffix, ipRangeEndSuffix int
	)
	ipRange := strings.Split(val, "-")
	if len(ipRange) == 2 {
		startIp = ipRange[0]
		startIpEndSuffix = StrToInt(strings.Split(startIp, ".")[3])
		ipRangeEndSuffix = StrToInt(ipRange[1])
	}
	return startIp, startIpEndSuffix, ipRangeEndSuffix, len(ipRange) == 2
}

func ExecutedTimeFormat(inSeconds int64) string {
	if inSeconds == 0 {
		return ""
	}
	minutes := inSeconds / 60
	seconds := inSeconds % 60
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}
