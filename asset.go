package offchain

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
