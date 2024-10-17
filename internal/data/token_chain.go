package data

type TokenChainsQ interface {
	New() TokenChainsQ
	Select() ([]TokenChain, error)
	Get() (*TokenChain, error)

	FilterByTokenID(ids ...string) TokenChainsQ
	FilterByChainID(ids ...string) TokenChainsQ
	FilterByBridgingType(types ...BridgingType) TokenChainsQ
}

type BridgingTypeQ interface {
	IsWrapped() bool
}

type TokenChain struct {
	TokenID         string
	ChainID         string       `fig:"chain_id,required"`
	ContractAddress *string      `fig:"contract_address"`
	TokenType       string       `fig:"token_type,required"`
	BridgingType    BridgingType `fig:"bridging_type,required"`
	AutoSend        bool         `fig:"auto_send"`
}

type BridgingType uint8

const (
	BridgingTypeLP BridgingType = iota
	BridgingTypeWrapped
	BridgingTypeUSDC
)

func (b BridgingType) IsLiquidPool() bool {
	if b > BridgingTypeUSDC {
		panic("unsupported bridging type")
	}
	return b == BridgingTypeLP
}

func (b BridgingType) IsWrapped() bool {
	if b > BridgingTypeUSDC {
		panic("unsupported bridging type")
	}
	return b == BridgingTypeWrapped
}

func (b BridgingType) IsUSDC() bool {
	if b > BridgingTypeUSDC {
		panic("unsupported bridging type")
	}
	return b == BridgingTypeUSDC
}
