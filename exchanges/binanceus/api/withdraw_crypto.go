package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type WithdrawRequest struct {
	Coin            oc.SymbolId  `json:"coin"`
	Network         oc.NetworkId `json:"network"`
	WithdrawOrderId string       `json:"withdrawOrderId,omitempty"`
	Address         oc.Address   `json:"address"`
	AddressTag      string       `json:"addressTag,omitempty"`
	Amount          oc.Amount    `json:"amount"`
	RecvWindow      *int64       `json:"recvWindow,omitempty"`
	TimestampMillis int64        `json:"timestamp,omitempty"`
}

type WithdrawResponse struct {
	Id string `json:"id"`
}

// https://docs.binance.us/#withdraw-crypto
func (c *Client) WithdrawCrypto(args *WithdrawRequest) (*WithdrawResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}

	var response WithdrawResponse
	query := url.Values{}
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))
	// Required parameters
	query.Set("coin", string(args.Coin))
	query.Set("network", string(args.Network))
	query.Set("address", string(args.Address))
	query.Set("amount", args.Amount.String())

	// Optional parameters
	if args.WithdrawOrderId != "" {
		query.Set("withdrawOrderId", args.WithdrawOrderId)
	}
	if args.AddressTag != "" {
		query.Set("addressTag", args.AddressTag)
	}
	if args.RecvWindow != nil {
		if *args.RecvWindow > 60000 {
			return nil, fmt.Errorf("recvWindow cannot be greater than 60000")
		}
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("POST", "/sapi/v1/capital/withdraw/apply", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
