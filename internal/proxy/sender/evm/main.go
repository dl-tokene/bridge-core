package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/signature"
	"math/big"
)

type TxSender struct {
	client  *ethclient.Client
	signer  signature.Signer
	chainID *big.Int
}

func NewTxSender(rpc string, signer signature.Signer) (Sender, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}

	return &TxSender{
		client:  client,
		signer:  signer,
		chainID: chainID,
	}, nil
}

func (s *TxSender) SendTx(tx *types.Transaction) (common.Hash, error) {
	tx, err := s.signer.SignTx(tx, s.chainID)
	if err != nil {
		return common.Hash{}, errors.Wrap(err, "failed to sign tx")
	}

	err = s.client.SendTransaction(context.TODO(), tx)

	return tx.Hash(), nil
}
