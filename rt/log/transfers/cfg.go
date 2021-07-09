package transfers

type TransferType int

const (
	CONSOLE TransferType = iota
	UDP
	TCP
	HTTP
)

type TransferFn func(ch chan []byte)

type TransferConfigure struct {
	Type   TransferType
	Server string
	Port   int
}
