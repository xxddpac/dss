package scan

import (
	"context"
	"dss/common/log"
	"dss/core/global"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

type SSH struct {
	WeakPasswordScan
}

func check(w WeakPasswordScan, ch chan struct{}) <-chan struct{} {
	go func() {
		client, err := ssh.Dial(global.TCP, fmt.Sprintf("%v:%v", w.Host, w.Port), &ssh.ClientConfig{
			User:            w.Username,
			Auth:            []ssh.AuthMethod{ssh.Password(w.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         2 * time.Second,
		})
		if err == nil {
			log.InfoF("found ssh weak password! host:%v,password:%v", w.Host, w.Password)
			_ = client.Close()
		}
		ch <- struct{}{}
	}()
	return ch
}

func (s *SSH) Do() {
	ch := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	/*
		Force disconnect if no response received within 5 seconds.
		fix ssh dial hang status through context.
		see: https://github.com/golang/go/issues/21941
	*/
	defer cancel()
	select {
	case <-ctx.Done():
		return
	case <-check(s.WeakPasswordScan, ch):
	}
}
