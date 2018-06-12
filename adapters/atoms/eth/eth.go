package eth

// EthereumData
type EthereumData struct {
	SwapID   [32]byte `json:"swap_id"`
	HashLock [32]byte `json:"hash_lock"`
}
