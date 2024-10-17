package signature

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gitlab.com/tokend/bridge/core/internal/data"
	"math/big"
	"testing"
)

func TestNativeLog(t *testing.T) {
	expectedResult := "0x6514d8c53feac8b0e102b3e25b8778fb849645e19e1badeadc2453b8232368a5"

	am, _ := big.NewInt(0).SetString("10000000000000000000", 10)
	log := NativeLog{
		Amount:     am,
		Receiver:   "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		TxHash:     common.HexToHash("0xc4f46c912cc2a1f30891552ac72871ab0f0e977886852bdd5dccd221a595647d"),
		EventIndex: 1794147,
		ChainID:    big.NewInt(31378),
	}

	hash := hexutil.Encode(log.Hash())

	if hash != expectedResult {
		t.Log(hash)
		t.Errorf("Wrong hash")
	}
}

func TestErc20Log(t *testing.T) {
	expectedResult := "0x9edbd32e49feb02e29e9c99a5a07b72b4d519e67d81cb3b360fd8d9a6ac4c267"

	am, _ := big.NewInt(0).SetString("100000000000000000000", 10)
	log := Erc20Log{
		TokenAddress: "0x5FbDB2315678afecb367f032d93F642f64180aa3",
		Amount:       am,
		Receiver:     "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		TxHash:       common.HexToHash("0xc4f46c912cc2a1f30891552ac72871ab0f0e977886852bdd5dccd221a595647d"),
		EventIndex:   1794147,
		ChainID:      big.NewInt(31378),
		BridgingType: data.BridgingTypeLP,
	}

	hash := hexutil.Encode(log.Hash())

	if hash != expectedResult {
		t.Log(hash)
		t.Errorf("Wrong hash")
	}
}

func TestErc721Log(t *testing.T) {
	expectedResult := "0xaed04cf9ba3e2ee69e3fa3adfa7fc3ffa8eb0d678847a9d20da397bef9c1d9de"

	log := Erc721Log{
		TokenAddress: "0x5FbDB2315678afecb367f032d93F642f64180aa3",
		TokenID:      big.NewInt(5000),
		Receiver:     "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		TxHash:       common.HexToHash("0xc4f46c912cc2a1f30891552ac72871ab0f0e977886852bdd5dccd221a595647d"),
		EventIndex:   1794147,
		ChainID:      big.NewInt(31378),
		TokenUri:     "https://some.link",
		BridgingType: data.BridgingTypeWrapped,
	}

	hash := hexutil.Encode(log.Hash())

	if hash != expectedResult {
		t.Log(hash)
		t.Errorf("Wrong hash")
	}
}

func TestErc1155Log(t *testing.T) {
	expectedResult := "0x21a2272dd7099ac1be823959374442ba4e04237196dc006419633481f4eeeaae"

	log := Erc1155Log{
		TokenAddress: "0x5FbDB2315678afecb367f032d93F642f64180aa3",
		TokenID:      big.NewInt(5000),
		Amount:       big.NewInt(10),
		Receiver:     "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		TxHash:       common.HexToHash("0xc4f46c912cc2a1f30891552ac72871ab0f0e977886852bdd5dccd221a595647d"),
		EventIndex:   1794147,
		ChainID:      big.NewInt(31378),
		TokenUri:     "https://some.link",
		BridgingType: data.BridgingTypeWrapped,
	}

	hash := hexutil.Encode(log.Hash())

	if hash != expectedResult {
		t.Log(hash)
		t.Errorf("Wrong hash")
	}
}
