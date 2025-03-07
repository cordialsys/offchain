package api

import (
	"net/url"

	oc "github.com/cordialsys/offchain"
)

type DepositAddressChain struct {
	ChainType         string       `json:"chainType"`
	AddressDeposit    string       `json:"addressDeposit"`
	TagDeposit        string       `json:"tagDeposit"`
	Chain             oc.NetworkId `json:"chain"`
	BatchReleaseLimit string       `json:"batchReleaseLimit"`
	ContractAddress   string       `json:"contractAddress"`
}

type GetDepositAddressResult struct {
	Coin   string                `json:"coin"`
	Chains []DepositAddressChain `json:"chains"`
}

type GetDepositAddressResponse = Response[GetDepositAddressResult]

// https://bybit-exchange.github.io/docs/v5/asset/deposit/master-deposit-addr
func (c *Client) GetMasterDepositAddress(coin oc.SymbolId, network oc.NetworkId) (*GetDepositAddressResponse, error) {
	params := url.Values{}
	params.Set("coin", string(coin))
	params.Set("chainType", string(network))
	var result GetDepositAddressResponse
	_, err := c.Request("GET", "/v5/asset/deposit/query-address", nil, &result, params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
