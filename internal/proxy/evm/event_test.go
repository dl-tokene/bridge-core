package evm_test

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"gitlab.com/tokend/bridge/core/internal/data"
	"gitlab.com/tokend/bridge/core/internal/data/mem"
	"gitlab.com/tokend/bridge/core/internal/ipfs"
	"gitlab.com/tokend/bridge/core/internal/proxy"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/signature"
	"gitlab.com/tokend/bridge/core/resources"
	"testing"
)

func setupEventTest() (*appConfig, error) {
	qTestnetRpc := "wss://rpc-ws.qtestnet.org"
	qTestnetTokenAddr := "0x340F12c1A1E30c00d502B9403328A09DD1153AbC"

	tokens := []data.Token{
		{
			ID:     "1",
			Name:   "USDC Testnet Token",
			Symbol: "USDC",
			Icon:   nil,
			Type:   resources.FUNGIBLE,
			Chains: []data.TokenChain{
				{
					ChainID:         "qtestnet",
					TokenType:       evm.TokenTypeErc20,
					BridgingType:    data.BridgingTypeLP,
					ContractAddress: &qTestnetTokenAddr,
					AutoSend:        false,
				},
			},
		},
	}

	chains := []data.Chain{
		{
			ID:             "qTestnet",
			Name:           "Q testnet",
			Icon:           nil,
			Type:           "evm",
			ChainParams:    nil,
			Confirmations:  12,
			BridgeContract: "0x0923BE8AC8F688382DF8B631978398C29b6e7474",
			RpcEndpoint:    qTestnetRpc,
			Tokens:         nil,
		},
		{
			ID:             "wrongQTestnet",
			Name:           "Q testnet",
			Icon:           nil,
			Type:           "evm",
			ChainParams:    nil,
			Confirmations:  12,
			BridgeContract: "",
			RpcEndpoint:    qTestnetRpc,
			Tokens:         nil,
		},
	}

	for i, chain := range chains {
		chain.Tokens = make([]data.TokenChain, 0)
		chains[i] = chain
	}

	tokenChains := make([]data.TokenChain, 0)
	for _, token := range tokens {
		for i, tokenChain := range token.Chains {
			tokenChain.TokenID = token.ID
			token.Chains[i] = tokenChain
			tokenChains = append(tokenChains, tokenChain)
			for k, chain := range chains {
				if chain.ID == tokenChain.ChainID {
					chain.Tokens = append(chain.Tokens, tokenChain)
					chains[k] = chain
				}
			}
		}
	}

	ipfs := ipfs.NewClient("https://ipfs.io")
	signer := "169e0067e34c4430b33113bf6806cced19b2655fd9e5b51c4847e31f5b171713"
	signerPk, err := crypto.HexToECDSA(signer)
	if err != nil {
		return nil, err
	}

	proxyRepo, err := proxy.NewProxyRepo(chains, signature.NewSigner(signerPk), ipfs)
	if err != nil {
		return nil, err
	}

	config := appConfig{
		ChainsQ:      mem.NewChainsQ(chains),
		TokensQ:      mem.NewTokenQ(tokens),
		TokenChainsQ: mem.NewTokenChainsQ(tokenChains),
		ProxyRepo:    proxyRepo,
		IPFS:         ipfs,
		Signer:       signature.NewSigner(signerPk),
	}

	return &config, nil
}

func TestGettingBridgeEvent(t *testing.T) {
	config, err := setupEventTest()
	require.NoError(t, err)

	proxy := config.ProxyRepo.Get("qTestnet")
	tokenChain, err := config.TokenChainsQ.
		FilterByTokenID("1").
		FilterByChainID("qtestnet").
		Get()
	require.NoError(t, err)

	txHash := "0x96ebed7275fca46e55a4062309038e8d3a838296346cc3d4dd85671932b50419"
	_, err = proxy.CheckFungibleLockEvent(txHash, 0, *tokenChain)
	require.NoError(t, err)
}

func TestGettingBridgeEventFail(t *testing.T) {
	config, err := setupEventTest()
	require.NoError(t, err)

	proxy := config.ProxyRepo.Get("wrongQTestnet")
	tokenChain, err := config.TokenChainsQ.
		FilterByTokenID("1").
		FilterByChainID("qtestnet").
		Get()
	require.NoError(t, err)

	txHash := "0x96ebed7275fca46e55a4062309038e8d3a838296346cc3d4dd85671932b50419"
	_, err = proxy.CheckFungibleLockEvent(txHash, 0, *tokenChain)
	require.Error(t, err)
}
