package http

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/republicprotocol/swapperd/foundation"
)

type Server interface {
	GetPing() GetPingResponse
	PostSwaps(PostSwapMessage) (PostSwapMessage, error)
}

type server struct {
	swapCh chan<- foundation.Swap
}

func NewServer(swapCh chan<- foundation.Swap) Server {
	return &server{
		swapCh: swapCh,
	}
}

func (server *server) GetPing() GetPingResponse {
	return GetPingResponse{
		Version: "0.1.0",
		SupportedTokens: []foundation.Token{
			foundation.TokenBTC,
			foundation.TokenETH,
			foundation.TokenWBTC,
		},
	}
}

func (server *server) PostSwaps(swapReq PostSwapMessage) (PostSwapMessage, error) {
	swap, err := decodePostSwap(swapReq)
	if err != nil {
		return PostSwapMessage{}, err
	}
	server.swapCh <- swap
	swapReq.SecretHash = MarshalSecretHash(swap.SecretHash)
	return swapReq, nil
}

func decodePostSwap(swap PostSwapMessage) (foundation.Swap, error) {
	secret := [32]byte{}
	if swap.ShouldInitiateFirst {
		rand.Read(secret[:])
		hash := sha256.Sum256(secret[:])
		swap.SecretHash = MarshalSecretHash(hash)
	}
	swapID, err := UnmarshalSwapID(swap.ID)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendToken, err := UnmarshalToken(swap.SendToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveToken, err := UnmarshalToken(swap.ReceiveToken)
	if err != nil {
		return foundation.Swap{}, nil
	}
	sendValue, err := UnmarshalAmount(swap.SendAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	receiveValue, err := UnmarshalAmount(swap.ReceiveAmount)
	if err != nil {
		return foundation.Swap{}, nil
	}
	secretHash, err := UnmarshalSecretHash(swap.SecretHash)
	if err != nil {
		return foundation.Swap{}, nil
	}
	return foundation.Swap{
		ID:                 swapID,
		Secret:             secret,
		SecretHash:         secretHash,
		TimeLock:           swap.TimeLock,
		SendToAddress:      swap.SendTo,
		ReceiveFromAddress: swap.ReceiveFrom,
		SendValue:          sendValue,
		ReceiveValue:       receiveValue,
		SendToken:          sendToken,
		ReceiveToken:       receiveToken,
		IsFirst:            swap.ShouldInitiateFirst,
	}, nil
}

// func (server *server) PostSwaps(swap PostSwap) (PostSwap, error) {
// }

// func (server *server) GetBalances() (Balances, error) {
// 	ethBal, err := ethereumBalance(
// 		server.config,
// 		server.keystr.GetKey(token.ETH).(keystore.EthereumKey),
// 	)
// 	if err != nil {
// 		return Balances{}, err
// 	}
// 	btcBal := bitcoinBalance(
// 		server.config,
// 		server.keystr.GetKey(token.BTC).(keystore.BitcoinKey),
// 	)
// 	return Balances{
// 		Ethereum: ethBal,
// 		Bitcoin:  btcBal,
// 	}, nil
// }

// func bitcoinBalance(conf config.Config, key keystore.BitcoinKey) Balance {
// 	conn := btc.NewConnWithConfig(conf.Bitcoin)
// 	balance := conn.Balance(key.AddressString, 0)
// 	return Balance{
// 		Address: key.AddressString,
// 		Amount:  strconv.FormatInt(balance, 10),
// 	}
// }

// func ethereumBalance(conf config.Config, key keystore.EthereumKey) (Balance, error) {
// 	conn, err := eth.NewConnWithConfig(conf.Ethereum)
// 	if err != nil {
// 		return Balance{}, err
// 	}
// 	bal, err := conn.Balance(key.Address)
// 	if err != nil {
// 		return Balance{}, err
// 	}
// 	return Balance{
// 		Address: key.Address.String(),
// 		Amount:  bal.String(),
// 	}, nil
// }
