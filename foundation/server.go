package foundation

type Balance struct {
	Address string `json:"address"`
	Amount  string `json:"balance"`
}

type Blockchain struct {
	Name    BlockchainName `json:"name"`
	Address string         `json:"address"`
}
