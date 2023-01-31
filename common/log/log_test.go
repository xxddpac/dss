package log

import (
	"testing"
)

var (
	fakeLogConfig = &Config{
		MaxSize:    10,
		Compress:   true,
		LogPath:    "",
		MaxAge:     0,
		MaxBackups: 0,
		LogLevel:   "info",
	}
)

func TestLogger(t *testing.T) {
	Init(fakeLogConfig)
	DebugF("debug")
	InfoF("info")
	WarnF("warn")
	Errorf("error")
}
