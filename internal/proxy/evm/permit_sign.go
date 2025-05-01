package evm

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm/generated/erc2612"
	"math/big"
	"time"
)

type TypeDefinition struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Eip712Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           string `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

type PermitMessage struct {
	Owner    string `json:"owner"`
	Spender  string `json:"spender"`
	Value    string `json:"value"`
	Nonce    string `json:"nonce"`
	Deadline string `json:"deadline"`
}

type PermitSign struct {
	Types       map[string][]TypeDefinition `json:"types"`
	PrimaryType string                      `json:"primaryType"`
	Domain      Eip712Domain                `json:"domain"`
	Message     PermitMessage               `json:"message"`
}

func (p *evmProxy) encodePermitSign(tokenAddress, owner string) (interface{}, error) {
	token, err := erc2612.NewErc2612(common.HexToAddress(tokenAddress), p.client)
	if err != nil {
		return nil, err
	}
	chainId, err := p.client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}
	tokenName, err := token.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	uin256max := big.NewInt(0)
	uin256max.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	nonce, err := token.Nonces(&bind.CallOpts{}, common.HexToAddress(owner))
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration("5m")
	if err != nil {
		return nil, err
	}
	deadline := big.NewInt(time.Now().Add(duration).Unix())

	types := map[string][]TypeDefinition{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "ChainId", Type: "uint256"},
			{Name: "VerifyingContract", Type: "address"},
		},
		"Permit": {
			{Name: "owner", Type: "address"},
			{Name: "spender", Type: "address"},
			{Name: "value", Type: "uint256"},
			{Name: "nonce", Type: "uint256"},
			{Name: "deadline", Type: "uint256"},
		},
	}

	domain := Eip712Domain{
		Name:              tokenName,
		Version:           "1",
		ChainID:           chainId.String(),
		VerifyingContract: tokenAddress,
	}

	message := PermitMessage{
		Owner:    owner,
		Spender:  p.bridgeContract.Hex(),
		Value:    uin256max.String(),
		Nonce:    nonce.String(),
		Deadline: deadline.String(),
	}

	permitSign := PermitSign{
		Types:       types,
		PrimaryType: "Permit",
		Domain:      domain,
		Message:     message,
	}

	return permitSign, nil
}
