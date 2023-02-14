package host

import (
	"context"
	"github.com/shirou/gopsutil/v3/host"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var (
	Name            atomic.Value
	PrivateIPv4     atomic.Value
	Platform        atomic.Value
	PlatformVersion atomic.Value
)

func init() {
	RefreshHost()
}

func LocalIP() string {
	PrivateIPv4S := PrivateIPv4.Load().([]string)
	if len(PrivateIPv4S) != 0 {
		return PrivateIPv4S[0]
	}
	return ""
}

func InitRefreshHost(ctx context.Context) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			RefreshHost()
		case <-ctx.Done():
			return
		}
	}
}

func RefreshHost() {
	var (
		privateIPv4               []string
		platform, platformVersion string
		address                   []net.Addr
		ip                        net.IP
	)
	hostname, _ := os.Hostname()
	Name.Store(hostname)
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, item := range interfaces {
		if strings.HasPrefix(item.Name, "docker") || strings.HasPrefix(item.Name, "lo") {
			continue
		}
		address, err = item.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range address {
			ip, _, err = net.ParseCIDR(addr.String())
			if err != nil || !ip.IsGlobalUnicast() {
				continue
			}
			if ip4 := ip.To4(); ip4 != nil {
				if (ip4[0] == 10) || (ip4[0] == 192 && ip4[1] == 168) || (ip4[0] == 172 && ip4[1]&0x10 == 0x10) {
					privateIPv4 = append(privateIPv4, ip4.String())
				}
			}
		}
	}
	if len(privateIPv4) > 5 {
		privateIPv4 = privateIPv4[:5]
	}
	PrivateIPv4.Store(privateIPv4)
	platform, _, platformVersion, _ = host.PlatformInformation()
	Platform.Store(platform)
	PlatformVersion.Store(platformVersion)
}
