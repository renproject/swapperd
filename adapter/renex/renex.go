package renex

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	bindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/adapter/keystore"
	"github.com/republicprotocol/renex-swapper-go/adapter/network"
	swapAdapter "github.com/republicprotocol/renex-swapper-go/adapter/swap"
	"github.com/republicprotocol/renex-swapper-go/adapter/watchdog"
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
	"github.com/republicprotocol/renex-swapper-go/service/logger"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/state"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type renexAdapter struct {
	state.State
	keystore.Keystore
	logger.Logger
	swap.Swapper
	network.Network
	*bindings.RenExSettlement
}

func New(config config.Config, keystore keystore.Keystore, network network.Network, watchdog watchdog.Watchdog, state state.State, logger logger.Logger) (renex.Adapter, error) {
	conn, err := ethclient.Connect(config.Ethereum)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Ethereum blockchain: %v", err)
	}

	renExSettlement, err := bindings.NewRenExSettlement(common.HexToAddress(config.RenEx.Settlement), bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, fmt.Errorf("cannot bind to RenEx accounts: %v", err)
	}

	return &renexAdapter{
		Keystore:        keystore,
		State:           state,
		Logger:          logger,
		Network:         network,
		RenExSettlement: renExSettlement,
		Swapper:         swap.NewSwapper(swapAdapter.New(config, keystore, network, watchdog, state, logger)),
	}, nil
}

// GetOrderMatch checks if a match is found and returns the match object. It
// keeps doing it until an order match is found or the waitTill time.
func (renex *renexAdapter) GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error) {
	for {
		PersonalOrder, ForeignOrder, ReceiveValue, SendValue, ReceiveCurrency, SendCurrency, err := renex.GetMatchDetails(&bind.CallOpts{}, orderID)
		if err != nil {
			return nil, err
		}

		if ReceiveCurrency != SendCurrency {
			return match.NewMatch(PersonalOrder, ForeignOrder, SendValue, ReceiveValue, token.Token(SendCurrency), token.Token(ReceiveCurrency)), nil
		}

		if time.Now().Unix() > waitTill {
			return nil, fmt.Errorf("Timed out")
		}
		time.Sleep(15 * time.Second)
	}
}

// GetAddress corresponding to the given token.
func (renex *renexAdapter) GetAddress(tok token.Token) []byte {
	switch tok {
	case token.BTC:
		return []byte(renex.Keystore.GetKey(tok).(keystore.BitcoinKey).AddressString)
	case token.ETH:
		return []byte(renex.Keystore.GetKey(tok).(keystore.EthereumKey).Address.String())
	default:
		panic("Unexpected token")
	}
}
