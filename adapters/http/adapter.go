package http

import "github.com/republicprotocol/atom-go/services/watch"

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

type PostOrder struct {
	OrderID   string `json:"orderID"`
	Signature string `json:"signature"`
}

type BoxHttpAdapter interface {
	WhoAmI(challenge string) (WhoAmI, error)
	PostOrder(order PostOrder) (PostOrder, error)
	BuildWatcher() (watch.Watch, error)
}
