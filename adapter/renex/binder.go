package renex

import (
	"fmt"
	"strings"
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
	fmt.Println("Waiting for the match to be found on RenEx")
	defer func() { fmt.Println("Match found") }()
	for {
		matchDetails, err := binder.GetMatchDetails(&bind.CallOpts{}, orderID)
		if err != nil || matchDetails.PriorityToken == matchDetails.SecondaryToken {
			if time.Now().Unix() > waitTill {
				return swap.Match{}, fmt.Errorf("Timed out")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		priorityToken, err := token.TokenCodeToToken(matchDetails.PriorityToken)
		if err != nil {
			return swap.Match{}, err
		}
		secondaryToken, err := token.TokenCodeToToken(matchDetails.SecondaryToken)
		if err != nil {
			return swap.Match{}, err
		}
		if matchDetails.OrderIsBuy {
			return swap.Match{
				PersonalOrderID: orderID,
				ForeignOrderID:  matchDetails.MatchedID,
				SendValue:       matchDetails.PriorityVolume.Add(matchDetails.PriorityVolume, matchDetails.PriorityFee),
				ReceiveValue:    matchDetails.SecondaryVolume.Add(matchDetails.SecondaryVolume, matchDetails.SecondaryFee),
				SendToken:       priorityToken,
				ReceiveToken:    secondaryToken,
			}, nil
		}
		return swap.Match{
			PersonalOrderID: orderID,
			ForeignOrderID:  matchDetails.MatchedID,
			SendValue:       matchDetails.SecondaryVolume.Add(matchDetails.SecondaryVolume, matchDetails.SecondaryFee),
			ReceiveValue:    matchDetails.PriorityVolume.Add(matchDetails.PriorityVolume, matchDetails.PriorityFee),
			SendToken:       secondaryToken,
			ReceiveToken:    priorityToken,
		}, nil
	}
}

func (binder *binder) verifyOrder(orderID [32]byte, waitTill int64) error {
	for {
		addr, err := binder.Orderbook.OrderTrader(&bind.CallOpts{}, orderID)
		if err != nil || addr.String() == "0x0000000000000000000000000000000000000000" {
			if time.Now().Unix() > waitTill {
				return fmt.Errorf("Timed out")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		for _, authorizedAddr := range binder.AuthorizedAddresses {
			if strings.ToLower(addr.String()) == strings.ToLower(authorizedAddr) {
				return nil
			}
			fmt.Printf("Expected submitting Trader Address %s to equal Authorized trader Address %s\n", addr.String(), authorizedAddr)
		}
	}
}
