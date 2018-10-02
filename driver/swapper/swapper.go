package swapper

type Swapper interface {
	Http(port int64)
	Withdraw(tk, to string, value, fee float64)
}
type swapper struct {
}

func NewSwapper() Swapper {
	return &swapper{}
}
