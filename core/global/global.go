package global

import "context"

const (
	PortScanQueue = "port_scan"
	Scan          = "scan"
	ScanRule      = "scan_rule"
	GrpcClient    = "grpc_client"
	SSH           = "22"
	REDIS         = "6379"
	MYSQL         = "3306"
	ROOT          = "root"
	TCP           = "tcp"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
