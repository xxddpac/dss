package buffer

import (
	"dss/core/models"
	"fmt"
	"testing"
	"time"
)

func TestBuffer(t *testing.T) {
	go func() {
		for i := 1; i <= 100; i++ {
			rec := models.Scan{
				Host: "1.1.1.1",
				Port: fmt.Sprintf("%d", i),
			}
			_ = WriteRecord(&rec, false)
			time.Sleep(time.Second)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 5)
			resp := ReadRecords()
			fmt.Println(len(resp))
			for _, item := range resp {
				fmt.Println(item.Host, item.Port)
			}
		}
	}()
	time.Sleep(12 * time.Second)
}
