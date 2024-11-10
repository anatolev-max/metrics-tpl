package enum

type MetricType int

const (
	Counter MetricType = iota + 1
	Gauge
)

func (m MetricType) String() string {
	return [...]string{
		"counter",
		"gauge",
	}[m-1]
}

func (m MetricType) EnumIndex() int {
	return int(m)
}
