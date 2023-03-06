package buffer

import (
	"dss/core/models"
	"errors"
	"sync"
)

var (
	mu                = &sync.Mutex{}
	buf               = [30000]*models.Scan{}
	offset            = 0
	ErrBufferOverflow = errors.New("buffer overflow")
	hook              func(interface{}) interface{}
)

func WriteRecord(rec *models.Scan, tolerate bool) (err error) {
	if hook != nil {
		rec = hook(rec).(*models.Scan)
	}
	mu.Lock()
	defer mu.Unlock()
	if offset < len(buf) {
		buf[offset] = rec
		offset++
		return
	}
	if tolerate {
		err = ErrBufferOverflow
	} else {
		buf[0] = rec
	}
	return
}

func ReadRecords() (ret []*models.Scan) {
	mu.Lock()
	defer mu.Unlock()
	ret = make([]*models.Scan, offset)
	copy(ret, buf[:offset])
	offset = 0
	return
}
