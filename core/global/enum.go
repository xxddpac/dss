package global

type RunMode int

const (
	Consumer RunMode = iota
	Producer
)
