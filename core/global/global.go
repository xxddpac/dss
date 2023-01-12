package global

import "context"

const (
	ScanQueue = "port_scan"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
