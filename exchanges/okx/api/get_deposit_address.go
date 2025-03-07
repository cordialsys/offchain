package api

import (
	"net/url"

	oc "github.com/cordialsys/offchain"
)

type DepositAddressData struct {
	// E.g.  "BTC-Bitcoin", "SOL-Solana"
	SymbolAndChain string `json:"chain"`

	ContractAddr string `json:"ctAddr"`
	Currency     string `json:"ccy"`
	To           string `json:"to"`
	Address      string `json:"addr"`
	VerifiedName string `json:"verifiedName"`
	Selected     bool   `json:"selected"`
}

type GetDepositAddressResponse = Response[[]DepositAddressData]

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-asset-bills-details
func (c *Client) GetDepositAddress(coin oc.SymbolId) (*GetDepositAddressResponse, error) {
	params := url.Values{}
	params.Set("ccy", string(coin))

	var response GetDepositAddressResponse
	_, err := c.Request("GET", "/api/v5/asset/deposit-address", nil, &response, params)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
