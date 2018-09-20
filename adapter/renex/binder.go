package renex

import (
	"bytes"
	"fmt"
	"time"

	"github.com/republicprotocol/renex-swapper-go/domain/match"
	"github.com/republicprotocol/renex-swapper-go/domain/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	bindings "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/bindings/eth"
	ethclient "github.com/republicprotocol/renex-swapper-go/adapter/blockchain/clients/eth"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/domain/order"
)

var (
	ErrVerificationFailed = fmt.Errorf("Given order id does not exist or belong to an authorized trader")
)

type Binder interface {
	GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error)
}

type binder struct {
	config.Config
	*bindings.Orderbook
	*bindings.RenExSettlement
}

func NewBinder(conf config.Config) (Binder, error) {
	conn, err := ethclient.NewConnWithConfig(conf.Ethereum)
	if err != nil {
		return nil, err
	}

	settlement, err := bindings.NewRenExSettlement(common.HexToAddress(conf.RenEx.Settlement), bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	orderbook, err := bindings.NewOrderbook(common.HexToAddress(conf.RenEx.Orderbook), bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	return &binder{
		Config:          conf,
		Orderbook:       orderbook,
		RenExSettlement: settlement,
	}, nil
}

// GetOrderMatch checks if a match is found and returns the match object. It
// keeps doing it until an order match is found or the waitTill time.
func (binder *binder) GetOrderMatch(orderID order.ID, waitTill int64) (match.Match, error) {
	if err := binder.verifyOrder(orderID, waitTill); err != nil {
		return nil, err
	}
	for {
		matchDetails, err := binder.GetMatchDetails(&bind.CallOpts{}, orderID)
		if err != nil || matchDetails.PriorityToken == matchDetails.SecondaryToken {
			if time.Now().Unix() > waitTill {
				return nil, fmt.Errorf("Timed out")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		if matchDetails.OrderIsBuy {
			return match.NewMatch(orderID, matchDetails.MatchedID, matchDetails.PriorityVolume, matchDetails.SecondaryVolume, token.Token(matchDetails.PriorityToken), token.Token(matchDetails.SecondaryToken)), nil
		}
		return match.NewMatch(orderID, matchDetails.MatchedID, matchDetails.SecondaryVolume, matchDetails.PriorityVolume, token.Token(matchDetails.SecondaryToken), token.Token(matchDetails.PriorityToken)), nil
	}
}

func (binder *binder) verifyOrder(orderID order.ID, waitTill int64) error {
	for {
		addr, err := binder.Orderbook.OrderTrader(&bind.CallOpts{}, orderID)
		if err != nil || bytes.Compare(addr.Bytes(), []byte{}) == 0 {
			if time.Now().Unix() > waitTill {
				return fmt.Errorf("Timed out")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		for _, authorizedAddr := range binder.AuthorizedAddresses {
			if addr.String() == authorizedAddr {
				return nil
			}
		}
		return ErrVerificationFailed
	}
}
