package enum

type ContentType int

const (
	ApplicationJSON ContentType = iota + 1
	TextPlain
)

func (c ContentType) String() string {
	return [...]string{
		"application/json",
		"text/plain",
	}[c-1]
}

func (c ContentType) EnumIndex() int {
	return int(c)
}
