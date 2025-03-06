package offchain

type ExchangeId string

var (
	Okx ExchangeId = "okx"
)

type SymbolId string
type NetworkId string
type ContractAddress string

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

type ExchangeConfig struct {
	ExchangeId ExchangeId
	ApiKey     string
	SecretKey  string
	Passphrase string
}

type BalanceDetail struct {
	// These are the normalized IDs (need upstream asset API to populate this)
	// AssetId string `json:"asset_id"`
	// ChainId string `json:"chain_id"`

	// These are the original exchange IDs
	SymbolId  SymbolId  `json:"symbol_id"`
	NetworkId NetworkId `json:"network_id"`

	// Amount, accounted for decimals
	Available Amount `json:"available"`
	Frozen    Amount `json:"frozen"`
}

type Client interface {
	GetAssets() ([]*Asset, error)
	GetBalances(args GetBalanceArgs) ([]*BalanceDetail, error)
}
