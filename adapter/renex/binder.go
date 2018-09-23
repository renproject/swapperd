package renex

import (
	"bytes"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/renex-swapper-go/adapter/config"
	"github.com/republicprotocol/renex-swapper-go/domain/swap"
	"github.com/republicprotocol/renex-swapper-go/domain/token"
)

var (
	ErrVerificationFailed = fmt.Errorf("Given order id does not exist or belong to an authorized trader")
)

type Binder interface {
	GetOrderMatch(orderID [32]byte, waitTill int64) (swap.Match, error)
}

type binder struct {
	config.Config
	*Orderbook
	*RenExSettlement
}

func NewBinder(conf config.Config) (Binder, error) {
	conn, err := NewConnWithConfig(conf)
	if err != nil {
		return nil, err
	}

	settlement, err := NewRenExSettlement(conn.RenExSettlement, bind.ContractBackend(conn.Client))
	if err != nil {
		return nil, err
	}

	orderbook, err := NewOrderbook(common.HexToAddress(conf.RenEx.Orderbook), bind.ContractBackend(conn.Client))
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
func (binder *binder) GetOrderMatch(orderID [32]byte, waitTill int64) (swap.Match, error) {
	if err := binder.verifyOrder(orderID, waitTill); err != nil {
		return swap.Match{}, err
	}
	for {
		matchDetails, err := binder.GetMatchDetails(&bind.CallOpts{}, orderID)
		if err != nil || matchDetails.PriorityToken == matchDetails.SecondaryToken {
			if time.Now().Unix() > waitTill {
				return swap.Match{}, fmt.Errorf("Timed out")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		if matchDetails.OrderIsBuy {
			return swap.Match{
				PersonalOrderID: orderID,
				ForeignOrderID:  matchDetails.MatchedID,
				SendValue:       matchDetails.PriorityVolume,
				ReceiveValue:    matchDetails.SecondaryVolume,
				SendToken:       token.Token(matchDetails.PriorityToken),
				ReceiveToken:    token.Token(matchDetails.SecondaryToken),
			}, nil
		}
		return swap.Match{
			PersonalOrderID: orderID,
			ForeignOrderID:  matchDetails.MatchedID,
			SendValue:       matchDetails.SecondaryVolume,
			ReceiveValue:    matchDetails.PriorityVolume,
			SendToken:       token.Token(matchDetails.SecondaryToken),
			ReceiveToken:    token.Token(matchDetails.PriorityToken),
		}, nil
	}
}

func (binder *binder) verifyOrder(orderID [32]byte, waitTill int64) error {
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
