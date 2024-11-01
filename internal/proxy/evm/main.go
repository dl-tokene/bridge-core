package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"gitlab.com/tokend/bridge/core/internal/ipfs"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/generated/bridge"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/signature"
	"gitlab.com/tokend/bridge/core/internal/proxy/sender/evm"
	"gitlab.com/tokend/bridge/core/internal/proxy/types"
	"math/big"
)

const (
	TokenTypeNative  = "native"
	TokenTypeErc20   = "erc20"
	TokenTypeErc721  = "erc721"
	TokenTypeErc1155 = "erc1155"
)

func NewProxy(rpc string, signer signature.Signer, bridgeContract string, ipfs ipfs.Client, confirmations int, sender evm.Sender) (types.Proxy, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	chainID, err := client.ChainID(context.TODO())
	if err != nil {
		return nil, err
	}

	b, err := bridge.NewBridge(common.HexToAddress(bridgeContract), client)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create bridge contract for address %s", bridgeContract)
	}

	return &evmProxy{
		client:         client,
		signer:         signer,
		chainID:        chainID,
		bridgeContract: common.HexToAddress(bridgeContract),
		bridge:         b,
		ipfsClient:     ipfs,
		confirmations:  confirmations,
		sender:         sender,
	}, nil
}

type evmProxy struct {
	client         *ethclient.Client
	signer         signature.Signer
	chainID        *big.Int
	bridgeContract common.Address
	bridge         *bridge.Bridge
	ipfsClient     ipfs.Client
	confirmations  int
	sender         evm.Sender
}
