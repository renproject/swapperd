package foundation

type Token struct {
	Name       string
	Blockchain string
}

var (
	TokenBTC  = Token{"BTC", "Bitcoin"}
	TokenWBTC = Token{"WBTC", "Ethereum"}
)
