package enum

type CounterType int

const (
	PollCounter CounterType = iota + 1
	RandomValue
)

func (c CounterType) String() string {
	return [...]string{
		"PollCounter",
		"RandomValue",
	}[c-1]
}

func (c CounterType) EnumIndex() int {
	return int(c)
}
