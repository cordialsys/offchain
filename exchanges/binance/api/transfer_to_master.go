package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type SubToMasterTransferRequest struct {
	Asset           oc.SymbolId `json:"asset"`
	Amount          oc.Amount   `json:"amount"`
	RecvWindow      *int64      `json:"recvWindow,omitempty"`
	TimestampMillis int64       `json:"timestamp,omitempty"`
}

type SubToMasterTransferResponse struct {
	TxnId string `json:"txnId"`
}

// SubToMasterTransfer transfers assets from a sub-account to the master account
// https://binance-docs.github.io/apidocs/spot/en/#transfer-to-master-for-sub-account
func (c *Client) SubToMasterTransfer(args *SubToMasterTransferRequest) (*SubToMasterTransferResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}

	var response SubToMasterTransferResponse
	query := url.Values{}

	// Required parameters
	query.Set("asset", string(args.Asset))
	query.Set("amount", args.Amount.String())
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	// Optional parameters
	if args.RecvWindow != nil {
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("POST", "/sapi/v1/sub-account/transfer/subToMaster", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
