package scan

import "fmt"

type SSH struct {
	WeakPasswordScan
}

func (s *SSH) Do() {
	fmt.Println("receive task from SSH", s)
}
