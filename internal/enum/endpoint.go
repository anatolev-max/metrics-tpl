package enum

type HttpEndpoint int

const (
    GetEndpoint HttpEndpoint = iota + 1
    UpdateEndpoint
)

func (h HttpEndpoint) String() string {
    return [...]string{
        "/value",
        "/update",
    }[h-1]
}

func (h HttpEndpoint) EnumIndex() int {
    return int(h)
}
