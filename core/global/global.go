package global

import "context"

const (
	PortScanQueue = "port_scan"
	Scan          = "scan"
	ScanRule      = "scan_rule"
	SSH           = "22"
	REDIS         = "6379"
	MYSQL         = "3306"
	ROOT          = "root"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
