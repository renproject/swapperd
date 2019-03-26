package erc20

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/renproject/libeth-go"
	"github.com/renproject/swapperd/core/wallet/swapper/immediate"
	"github.com/renproject/swapperd/foundation/blockchain"
	"github.com/renproject/swapperd/foundation/swap"
	"github.com/renproject/tokens"
	"github.com/sirupsen/logrus"
)

type erc20SwapContractBinder struct {
	id             [32]byte
	account        libeth.Account
	speed          libeth.TxExecutionSpeed
	swap           swap.Swap
	logger         logrus.FieldLogger
	swapperAddress common.Address
	swapperBinder  *ERC20SwapContract
	erc20          libeth.ERC20
	cost           blockchain.Cost
}

// NewERC20SwapContractBinder returns a new ERC20 Atom instance
func NewERC20SwapContractBinder(account libeth.Account, swap swap.Swap, cost blockchain.Cost, logger logrus.FieldLogger) (immediate.Contract, error) {
	erc20, err := account.NewERC20(string(swap.Token.Name))
	if err != nil {
		return nil, err
	}

	swapperAddress, err := account.ReadAddress(fmt.Sprintf("%sSwap", swap.Token.Name))
	if err != nil {
		return nil, err
	}

	swapperBinder, err := NewERC20SwapContract(swapperAddress, bind.ContractBackend(account.EthClient()))
	if err != nil {
		return nil, err
	}

	id, err := swapperBinder.SwapID(&bind.CallOpts{}, swap.SecretHash, big.NewInt(swap.TimeLock))
	if err != nil {
		return nil, err
	}

	fields := logrus.Fields{}
	fields["SwapID"] = swap.ID
	fields["ContractID"] = base64.StdEncoding.EncodeToString(id[:])
	fields["Token"] = swap.Token.Name
	logger = logger.WithFields(fields)

	if _, ok := cost[tokens.NameETH]; !ok {
		cost[tokens.NameETH] = big.NewInt(0)
	}

	if _, ok := cost[swap.Token.Name]; !ok {
		cost[swap.Token.Name] = big.NewInt(0)
	}

	swap.Value = new(big.Int).Add(swap.Value, swap.BrokerFee)

	return &erc20SwapContractBinder{
		account:        account,
		swapperAddress: swapperAddress,
		swapperBinder:  swapperBinder,
		erc20:          erc20,
		speed:          libeth.TxExecutionSpeed(swap.Speed),
		logger:         logger,
		swap:           swap,
		id:             id,
		cost:           cost,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Initiate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		return err
	}

	if !initiatable {
		atom.logger.Info(fmt.Sprintf("Skipping initiate on Ethereum blockchain"))
		return nil
	}
	atom.logger.Info(fmt.Sprintf("Initiating on Ethereum blockchain"))

	// Approve the contract to transfer tokens
	approveTx, err := atom.erc20.Approve(ctx, atom.swapperAddress, atom.sendValue(), atom.speed)
	if err != nil {
		return err
	}
	atom.cost[tokens.NameETH] = new(big.Int).Add(atom.cost[tokens.NameETH], approveTx.Cost())

	// Initiate the Atomic Swap
	initiateTx, err := atom.account.Transact(
		ctx,
		atom.speed,
		nil,
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			var tx *types.Transaction
			var err error
			if atom.swap.BrokerFee.Cmp(big.NewInt(0)) > 0 {
				tx, err = atom.swapperBinder.InitiateWithFees(tops, atom.id, common.HexToAddress(atom.swap.SpendingAddress), common.HexToAddress(atom.swap.BrokerAddress), atom.swap.BrokerFee, atom.swap.SecretHash, big.NewInt(atom.swap.TimeLock), atom.sendValue())
				if err != nil {
					return tx, err
				}
				atom.cost[atom.swap.Token.Name] = new(big.Int).Add(atom.cost[atom.swap.Token.Name], atom.swap.BrokerFee)
			} else {
				tx, err = atom.swapperBinder.Initiate(tops, atom.id, common.HexToAddress(atom.swap.SpendingAddress), atom.swap.SecretHash, big.NewInt(atom.swap.TimeLock), atom.sendValue())
				if err != nil {
					return tx, err
				}
			}
			msg, _ := atom.account.FormatTransactionView("Initiated on Ethereum blockchain", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !initiatable
		},
		1,
	)
	if err != nil {
		if err == libeth.ErrPreConditionCheckFailed {
			return immediate.ErrAlreadyInitiated
		}
		return err
	}
	atom.cost[tokens.NameETH] = new(big.Int).Add(atom.cost[tokens.NameETH], initiateTx.Cost())
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Refund() error {
	atom.logger.Info("Refunding on Ethereum blockchain")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	tx, err := atom.account.Transact(
		ctx,
		atom.speed,
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return refundable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.swapperBinder.Refund(tops, atom.id)
			if err != nil {
				return nil, err
			}
			msg, _ := atom.account.FormatTransactionView("Refunded on Ethereum blockchain", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	)
	if err != nil && err != libeth.ErrPreConditionCheckFailed {
		return err
	}
	atom.cost[tokens.NameETH] = new(big.Int).Add(atom.cost[tokens.NameETH], tx.Cost())
	if _, ok := atom.cost[atom.swap.Token.Name]; ok {
		atom.cost[atom.swap.Token.Name] = new(big.Int).Sub(atom.cost[atom.swap.Token.Name], atom.swap.BrokerFee)
	}
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) AuditSecret() ([32]byte, error) {
	atom.logger.Info("Auditing secret on Ethereum blockchain")
	redeemable, err := atom.swapperBinder.Redeemable(&bind.CallOpts{}, atom.id)
	if err != nil {
		atom.logger.Error(err)
		return [32]byte{}, err
	}
	if redeemable {
		if time.Now().Unix() > atom.swap.TimeLock {
			atom.logger.Error(immediate.ErrSwapExpired)
			return [32]byte{}, immediate.ErrSwapExpired
		}
		return [32]byte{}, immediate.ErrAuditPending
	}

	secret, err := atom.swapperBinder.AuditSecret(&bind.CallOpts{}, atom.id)
	if err != nil {
		return [32]byte{}, err
	}

	atom.logger.Info(fmt.Sprintf("Audit succeeded on Ethereum blockchain secret = %s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Audit() error {
	atom.logger.Info(fmt.Sprintf("Waiting for initiation on Ethereum blockchain"))
	initiatable, err := atom.swapperBinder.Initiatable(&bind.CallOpts{}, atom.id)
	if err != nil {
		atom.logger.Error(err)
		return err
	}

	if initiatable {
		if time.Now().Unix() > atom.swap.TimeLock {
			atom.logger.Error(immediate.ErrSwapExpired)
			return immediate.ErrSwapExpired
		}
		return immediate.ErrAuditPending
	}

	auditReport, err := atom.swapperBinder.Audit(&bind.CallOpts{}, atom.id)
	if err != nil {
		atom.logger.Error(err)
		return err
	}

	if auditReport.To.String() != atom.swap.SpendingAddress {
		err := fmt.Errorf("Receiver Address Mismatch Expected: %v Actual: %v", atom.swap.SpendingAddress, auditReport.To.String())
		atom.logger.Error(err)
		return err
	}

	value := new(big.Int).Sub(atom.swap.Value, atom.swap.BrokerFee)
	if auditReport.Value.Cmp(value) < 0 {
		return fmt.Errorf("Receive value mismatch: expected %v, got %v", atom.swap.Value, auditReport.Value)
	}
	atom.logger.Info(fmt.Sprintf("Audit successful on Ethereum blockchain"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *erc20SwapContractBinder) Redeem(secret [32]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	tx, err := atom.account.Transact(
		ctx,
		atom.speed,
		func() bool {
			redeemable, err := atom.swapperBinder.Redeemable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return redeemable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := atom.swapperBinder.Redeem(tops, atom.id, common.HexToAddress(atom.swap.WithdrawAddress), secret)
			if err != nil {
				return nil, err
			}

			msg, _ := atom.account.FormatTransactionView("Redeemed the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.swapperBinder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	)
	if err != nil {
		if err != libeth.ErrPreConditionCheckFailed {
			return err
		}
		atom.logger.Info("Skipping redeem on Ethereum Blockchain")
	}
	atom.cost[tokens.NameETH] = new(big.Int).Add(atom.cost[tokens.NameETH], tx.Cost())
	return nil
}

func (atom *erc20SwapContractBinder) Cost() blockchain.Cost {
	return atom.cost
}

func (atom *erc20SwapContractBinder) sendValue() *big.Int {
	if additionalFee := atom.swap.Token.AdditionalTransactionFee(atom.swap.Value); additionalFee != nil {
		return new(big.Int).Add(atom.swap.Value, additionalFee)
	}
	return atom.swap.Value
}
