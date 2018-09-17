package renex

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	bindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	swapAdapter "github.com/republicprotocol/renex-swapper-go/adapter/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
	"github.com/republicprotocol/renex-swapper-go/service/renex"
	"github.com/republicprotocol/renex-swapper-go/service/swap"
)

type renexAdapter struct {
	swap.Builder
	*bindings.RenExSettlement
}

func New(config config.Config, network swap.Network, watchdog swap.Watchdog, state renex.State, logger swap.Logger, atomBuilder swap.AtomBuilder) (renex.Adapter, error) {
	conn, err := ethclient.Connect(config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to Ethereum blockchain: %v", err)
	}

	renExSettlement, err := bindings.NewRenExSettlement(conn.RenExSettlement, bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, fmt.Errorf("cannot bind to RenEx accounts: %v", err)
	}

	return &renexAdapter{
		RenExSettlement: renExSettlement,
		Builder:         swapAdapter.New(network, watchdog, state, logger, atomBuilder),
	}, nil
}

func (renex *renexAdapter) NewSwapAdapter(req swap.Request) (swap.Adapter, error) {
	return renex.Builder.New(req)
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
			return match.NewMatch(PersonalOrder, ForeignOrder, SendValue, ReceiveValue, SendCurrency, ReceiveCurrency), nil
		}

		if time.Now().Unix() > waitTill {
			return nil, fmt.Errorf("Timed out")
		}

		time.Sleep(15 * time.Second)
	}
}
