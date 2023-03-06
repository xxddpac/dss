package scan

import (
	"dss/common/log"
	"dss/core/dao"
	"dss/core/global"
	"encoding/json"
	"go.uber.org/zap"
	"time"
)

type WeakPasswordScan struct {
	scanInfo
	Username string
	Password string
}

func Dispatch() {
	log.Info("start run dispatch...")
	for {
		if msg, err := dao.Redis.BRPop(global.IpScan, 0*time.Second); err == nil {
			w := &WeakPasswordScan{}
			if err := json.Unmarshal([]byte(msg), w); err != nil {
				log.Error("failed unmarshal json at parse message", zap.String("msg", msg), zap.Error(err))
				continue
			}
			switch w.Port {
			case global.SSH:
				//poolForDispatch.Add(&SSH{*w})
			case global.MYSQL:

			case global.REDIS:

			default:

			}
		}
	}
}
