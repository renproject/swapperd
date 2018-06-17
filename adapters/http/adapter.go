package http

type BoxInfo struct {
	challenge           string   `json:"challenge"`
	version             string   `json:"version"`
	supportedCurrencies []string `json:"supportedCurrencies"`
}

type WhoAmI struct {
	boxInfo   BoxInfo `json:"boxInfo"`
	signature string  `json:"signature"`
}

type PostOrder struct {
	OrderID   string `json:"orderID"`
	Signature string `json:"signature"`
}

type BoxHttpAdapter interface {
	WhoAmI(challenge string) (WhoAmI, error)
	PostOrder(order PostOrder) (PostOrder, error)
}
