package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Sender interface {
	SendTx(tx *types.Transaction) (common.Hash, error)
}
