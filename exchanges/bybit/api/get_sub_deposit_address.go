package api

import (
	"net/url"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
)

// https://bybit-exchange.github.io/docs/v5/asset/deposit/sub-deposit-addr
func (c *Client) GetSubDepositAddress(accountId client.AccountId, coin oc.SymbolId, network oc.NetworkId) (*GetDepositAddressResponse, error) {
	params := url.Values{}
	params.Set("coin", string(coin))
	params.Set("chainType", string(network))
	params.Set("subMemberId	", string(accountId))
	var result GetDepositAddressResponse
	_, err := c.Request("GET", "/v5/asset/deposit/query-sub-member-address", nil, &result, params)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
