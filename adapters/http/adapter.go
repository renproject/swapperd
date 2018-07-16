package http

type BoxInfo struct {
	Challenge           string   `json:"challenge"`
	Version             string   `json:"version"`
	AuthorizedAddresses []string `json:"authorizedAddresses"`
	SupportedCurrencies []string `json:"supportedCurrencies"`
}

type WhoAmI struct {
	BoxInfo   BoxInfo `json:"boxInfo"`
	Signature string  `json:"signature"`
}

type Status struct {
	OrderID string `json:"orderID"`
	Status  string `json:"status"`
}

type PostOrder struct {
	OrderID   string `json:"orderID"`
	Signature string `json:"signature"`
}

type Balance struct {
	Address string `json:"address"`
	Amount  uint64 `json:"amount"`
}

type Balances []Balance

type BoxHttpAdapter interface {
	WhoAmI(challenge string) (WhoAmI, error)
	PostOrder(order PostOrder) (PostOrder, error)
	GetStatus(orderID string) (Status, error)
	GetBalances() (Balances, error)
}
