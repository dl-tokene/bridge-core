package types

import (
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

type Sender interface {
	SendTx(tx *ethTypes.Transaction, chainId *big.Int) (interface{}, error)
}
