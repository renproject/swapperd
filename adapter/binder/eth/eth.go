package eth

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/beth-go"
	"github.com/republicprotocol/swapperd/core/swapper/immediate"
	"github.com/republicprotocol/swapperd/foundation/blockchain"
	"github.com/republicprotocol/swapperd/foundation/swap"
	"github.com/sirupsen/logrus"
)

type ethSwapContractBinder struct {
	id      [32]byte
	account beth.Account
	swap    swap.Swap
	logger  logrus.FieldLogger
	binder  *EthSwapContract
	cost    blockchain.Cost
}

// NewETHSwapContractBinder returns a new Ethereum RequestAtom instance
func NewETHSwapContractBinder(account beth.Account, swap swap.Swap, cost blockchain.Cost, logger logrus.FieldLogger) (immediate.Contract, error) {
	swapperAddr, err := account.ReadAddress("ETHSwapContract")
	if err != nil {
		return nil, err
	}

	contract, err := NewEthSwapContract(swapperAddr, bind.ContractBackend(account.EthClient()))
	if err != nil {
		return nil, err
	}

	id, err := contract.SwapID(&bind.CallOpts{}, swap.SecretHash, big.NewInt(swap.TimeLock))
	if err != nil {
		return nil, err
	}

	fields := logrus.Fields{}
	fields["SwapID"] = swap.ID
	fields["ContractID"] = base64.StdEncoding.EncodeToString(id[:])
	fields["Token"] = swap.Token.Name
	logger = logger.WithFields(fields)

	if _, ok := cost[blockchain.ETH]; !ok {
		cost[blockchain.ETH] = big.NewInt(0)
	}

	logger.Info(swap.ID, fmt.Sprintf("Ethereum Atomic Swap ID: %s", base64.StdEncoding.EncodeToString(id[:])))
	return &ethSwapContractBinder{
		account: account,
		binder:  contract,
		logger:  logger,
		swap:    swap,
		id:      id,
		cost:    cost,
	}, nil
}

// Initiate a new Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Initiate() error {
	atom.logger.Info("Initiating")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initiate the Atomic Swap
	if err := atom.account.Transact(
		ctx,
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return initiatable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tops.GasPrice = atom.swap.Fee
			tops.Value = atom.swap.Value
			var tx *types.Transaction
			var err error
			if atom.swap.BrokerFee.Cmp(big.NewInt(0)) > 0 {
				tx, err = atom.binder.InitiateWithFees(tops, atom.id, common.HexToAddress(atom.swap.SpendingAddress), common.HexToAddress(atom.swap.BrokerAddress), atom.swap.BrokerFee, atom.swap.SecretHash, big.NewInt(atom.swap.TimeLock), atom.swap.Value)
				if err != nil {
					return tx, err
				}

				atom.cost[blockchain.ETH] = new(big.Int).Add(atom.cost[blockchain.ETH], atom.swap.BrokerFee)
			} else {
				tx, err = atom.binder.Initiate(tops, atom.id, common.HexToAddress(atom.swap.SpendingAddress), atom.swap.SecretHash, big.NewInt(atom.swap.TimeLock), atom.swap.Value)
				if err != nil {
					return tx, err
				}
			}

			txFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
			atom.cost[blockchain.ETH] = new(big.Int).Add(atom.cost[blockchain.ETH], txFee)

			tops.Value = big.NewInt(0)
			msg, _ := atom.account.FormatTransactionView("Initiated the atomic swap", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !initiatable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

// Refund an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Refund() error {
	atom.logger.Info("Refunding the atomic swap")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := atom.account.Transact(
		ctx,
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return refundable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tops.GasPrice = atom.swap.Fee
			tx, err := atom.binder.Refund(tops, atom.id)
			if err != nil {
				return nil, err
			}

			txFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
			atom.cost[blockchain.ETH] = new(big.Int).Add(atom.cost[blockchain.ETH], txFee)

			msg, _ := atom.account.FormatTransactionView("Refunded the atomic swap", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

// AuditSecret audits the secret of an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) AuditSecret() ([32]byte, error) {
	atom.logger.Info("Auditing secret on ethereum blockchain")
	redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
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

	secret, err := atom.binder.AuditSecret(&bind.CallOpts{}, atom.id)
	if err != nil {
		return [32]byte{}, err
	}
	atom.logger.Info(fmt.Sprintf("Audit success on ethereum blockchain secret=%s", base64.StdEncoding.EncodeToString(secret[:])))
	return secret, nil
}

// Audit an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Audit() error {
	atom.logger.Info(fmt.Sprintf("Waiting for initiation on ethereum blockchain"))
	initiatable, err := atom.binder.Initiatable(&bind.CallOpts{}, atom.id)
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
	auditReport, err := atom.binder.Audit(&bind.CallOpts{}, atom.id)
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
	if auditReport.Value.Cmp(value) != 0 {
		atom.logger.Error(fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.swap.Value, auditReport.Value))
		return fmt.Errorf("Receive Value Mismatch Expected: %v Actual: %v", atom.swap.Value, auditReport.Value)
	}
	atom.logger.Info(fmt.Sprintf("Audit successful"))
	return nil
}

// Redeem an Atom swap by calling a function on ethereum
func (atom *ethSwapContractBinder) Redeem(secret [32]byte) error {
	atom.logger.Info("Redeeming the atomic swap")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := atom.account.Transact(
		ctx,
		func() bool {
			redeemable, err := atom.binder.Redeemable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return redeemable
		},
		func(tops *bind.TransactOpts) (*types.Transaction, error) {
			tops.GasPrice = atom.swap.Fee
			tx, err := atom.binder.Redeem(tops, atom.id, common.HexToAddress(atom.swap.WithdrawAddress), secret)
			if err != nil {
				return nil, err
			}

			txFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
			atom.cost[blockchain.ETH] = new(big.Int).Add(atom.cost[blockchain.ETH], txFee)

			msg, _ := atom.account.FormatTransactionView("Redeemed the atomic swap on Ethereum blockchain", tx.Hash().String())
			atom.logger.Info(msg)
			return tx, nil
		},
		func() bool {
			refundable, err := atom.binder.Refundable(&bind.CallOpts{}, atom.id)
			if err != nil {
				return false
			}
			return !refundable
		},
		1,
	); err != nil && err != beth.ErrPreConditionCheckFailed {
		return err
	}
	return nil
}

func (atom *ethSwapContractBinder) Cost() blockchain.Cost {
	return atom.cost
}
