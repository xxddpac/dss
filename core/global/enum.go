package global

type RunMode int

const (
	Consumer RunMode = iota
	Producer
)

type RuleType int

const (
	Single RuleType = iota + 1
	Range
	Cidr
)

func (r RuleType) String() string {
	switch r {
	case Single:
		return "单个IP类型" //192.168.1.1
	case Range:
		return "连续范围IP类型" //192.168.1.1-100
	case Cidr:
		return "网段类型" //192.168.1.0/24
	default:
		return "未知"
	}
}

type Level int

const (
	Critical Level = iota + 1
	High
	Middle
	Low
)

func (l Level) String() string {
	switch l {
	case Critical:
		return "严重"
	case High:
		return "高危"
	case Middle:
		return "中危"
	case Low:
		return "低危"
	default:
		return "未知"
	}
}
