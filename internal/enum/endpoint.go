package enum

type HTTPEndpoint int

const (
	GetEndpoint HTTPEndpoint = iota + 1
	UpdateEndpoint
)

func (h HTTPEndpoint) String() string {
	return [...]string{
		"/value",
		"/update",
	}[h-1]
}

func (h HTTPEndpoint) EnumIndex() int {
	return int(h)
}
