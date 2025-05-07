package sending

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/signature"
	"gitlab.com/tokend/bridge/core/internal/proxy/types"
	"gitlab.com/tokend/bridge/core/resources"
	"math/big"
)

func NewSender(client *ethclient.Client, signer signature.Signer) (types.Sender, error) {
	return &evmSender{
		client: client,
		signer: signer,
	}, nil
}

type evmSender struct {
	client *ethclient.Client
	signer signature.Signer
}

func (s *evmSender) SendTx(tx *ethTypes.Transaction, chainId *big.Int) (interface{}, error) {
	tx, err := s.signer.SignTx(tx, chainId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign transaction")
	}

	err = s.client.SendTransaction(context.TODO(), tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send transaction")
	}

	return encodeSentTx(tx.Hash(), chainId.String()), nil
}

func encodeSentTx(txHash common.Hash, chain string) interface{} {
	return resources.ProcessedTransaction{
		Key: resources.Key{
			ID:   txHash.String(),
			Type: resources.PROCESSED_TRANSACTION,
		},
		Relationships: resources.ProcessedTransactionRelationships{
			Chain: resources.Key{
				ID:   chain,
				Type: resources.CHAIN,
			}.AsRelation(),
		},
	}
}
