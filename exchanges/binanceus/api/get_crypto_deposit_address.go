package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type GetDepositAddressRequest struct {
	Coin            oc.SymbolId  `json:"coin"`
	Network         oc.NetworkId `json:"network,omitempty"`
	RecvWindow      *int64       `json:"recvWindow,omitempty"`
	TimestampMillis int64        `json:"timestamp,omitempty"`
}

type DepositAddressResponse struct {
	Coin    oc.SymbolId `json:"coin"`
	Address oc.Address  `json:"address"`
	Tag     string      `json:"tag"`
	URL     string      `json:"url"`
}

// https://docs.binance.us/#get-crypto-deposit-address
func (c *Client) GetDepositAddress(args *GetDepositAddressRequest) (*DepositAddressResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}

	var response DepositAddressResponse
	query := url.Values{}
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	// Required parameter
	query.Set("coin", string(args.Coin))

	// Optional parameters
	if args.Network != "" {
		query.Set("network", string(args.Network))
	}

	if args.RecvWindow != nil {
		if *args.RecvWindow > 60000 {
			return nil, fmt.Errorf("recvWindow cannot be greater than 60000")
		}
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("GET", "/sapi/v1/capital/deposit/address", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
