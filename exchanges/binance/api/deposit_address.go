package api

import (
	"fmt"
	"net/url"
	"time"

	oc "github.com/cordialsys/offchain"
)

type DepositAddressRequest struct {
	Coin    oc.SymbolId  `json:"coin"`
	Network oc.NetworkId `json:"network,omitempty"`
	Amount  *oc.Amount   `json:"amount,omitempty"`
}

type DepositAddressResponse struct {
	Address oc.Address  `json:"address"`
	Coin    oc.SymbolId `json:"coin"`
	Tag     string      `json:"tag"`
	Url     string      `json:"url"`
}

// https://developers.binance.com/docs/wallet/capital/deposite-address
func (c *Client) GetDepositAddress(args *DepositAddressRequest) (*DepositAddressResponse, error) {
	var response DepositAddressResponse
	query := url.Values{}

	// Required parameter
	query.Set("coin", string(args.Coin))
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

	// Optional parameters
	if args.Network != "" {
		query.Set("network", string(args.Network))
	}
	if args.Amount != nil {
		query.Set("amount", args.Amount.String())
	}

	_, err := c.Request("GET", "/sapi/v1/capital/deposit/address", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
