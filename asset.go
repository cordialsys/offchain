package offchain

type ExchangeId string

var (
	Okx     ExchangeId = "okx"
	Binance ExchangeId = "binance"
	Bybit   ExchangeId = "bybit"
)

var ValidExchangeIds = []ExchangeId{Okx, Binance, Bybit}

type SymbolId string
type NetworkId string
type ContractAddress string
type Address string

type Asset struct {
	SymbolId        SymbolId        `json:"symbol_id"`
	NetworkId       NetworkId       `json:"network_id"`
	ContractAddress ContractAddress `json:"contract_address"`
}

func NewAsset(symbolId SymbolId, networkId NetworkId, contractAddress ContractAddress) *Asset {
	return &Asset{
		SymbolId:        symbolId,
		NetworkId:       networkId,
		ContractAddress: contractAddress,
	}
}

type BalanceDetail struct {
	// These are the normalized IDs (need upstream asset API to populate this)
	// AssetId string `json:"asset_id"`
	// ChainId string `json:"chain_id"`

	// These are the original exchange IDs
	SymbolId SymbolId `json:"symbol_id"`
	// often the network is not relevant, as assets can be withdrawn to multiple networks
	NetworkId NetworkId `json:"network_id,omitempty"`

	// Amount, accounted for decimals
	Available Amount `json:"available"`
	Frozen    Amount `json:"frozen"`
}
