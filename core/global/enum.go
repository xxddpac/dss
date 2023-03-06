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

type TaskStatus int

const (
	Waiting TaskStatus = iota + 1
	Running
	Finished
	Error
)

func (t TaskStatus) String() string {
	switch t {
	case Waiting:
		return "等待中"
	case Running:
		return "检测中"
	case Finished:
		return "已完成"
	case Error:
		return "出错"
	default:
		return "未知"
	}
}
