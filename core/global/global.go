package global

import "context"

const (
	TimeUnix      = 100
	PortScanQueue = "port_scan"
	IpScan        = "ip_scan"
	Scan          = "scan"
	ScanRule      = "scan_rule"
	ScanTask      = "scan_task"
	GrpcClient    = "grpc_client"
	SSH           = "22"
	REDIS         = "6379"
	MYSQL         = "3306"
	TCP           = "tcp"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
