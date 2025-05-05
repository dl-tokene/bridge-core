package evm

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/internal/data"
)

const NativeTokenDecimals int = 18

func (p *evmProxy) Decimals(tokenChain data.TokenChain) (int, error) {
	switch tokenChain.TokenType {
	case TokenTypeNative:
		return NativeTokenDecimals, nil
	case TokenTypeErc20:
		return p.getDecimals(*tokenChain.ContractAddress)
	default:
		return 0, errors.Errorf("unsupported token type: %s, token: %s", tokenChain.TokenType, tokenChain.TokenID)
	}
}
