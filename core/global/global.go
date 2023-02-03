package global

import "context"

const (
	ScanQueue    = "port_scan"
	PortScan     = "port_scan"
	PortScanRule = "port_scan_rule"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
