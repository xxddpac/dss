package pprof

import (
	"dss/common/log"
	"fmt"
	"net/http"
)

func Pprof(enable bool, port int) {
	if enable {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Errorf("pprof ListenAndServe Error %s", err.Error())
		}
	}
}
