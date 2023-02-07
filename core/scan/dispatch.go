package scan

import (
	"dss/common/wp"
	"dss/core/global"
)

type WeakPasswordScan struct {
	scanInfo
	Username string
	Password string
}

func dispatch(scanInfo scanInfo) {
	for _, password := range wp.WeakPasswordList {
		w := WeakPasswordScan{scanInfo, global.ROOT, password}
		switch scanInfo.Port {
		case global.SSH:
			poolForPortScan.Add(&SSH{w})
		case global.MYSQL:
		case global.REDIS:
		}
	}
}
