package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type SubToSubTransferRequest struct {
	ToEmail         string      `json:"toEmail"`
	Asset           oc.SymbolId `json:"asset"`
	Amount          oc.Amount   `json:"amount"`
	RecvWindow      *int64      `json:"recvWindow,omitempty"`
	TimestampMillis int64       `json:"timestamp,omitempty"`
}

type SubToSubTransferResponse struct {
	TxnId string `json:"txnId"`
}

// https://developers.binance.com/docs/sub_account/asset-management/Transfer-to-Master
func (c *Client) SubToSubTransfer(args *SubToSubTransferRequest) (*SubToSubTransferResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}

	var response SubToSubTransferResponse
	query := url.Values{}

	// Required parameters
	query.Set("toEmail", args.ToEmail)
	query.Set("asset", string(args.Asset))
	query.Set("amount", args.Amount.String())
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	// Optional parameters
	if args.RecvWindow != nil {
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("POST", "/sapi/v1/sub-account/transfer/subToSub", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
