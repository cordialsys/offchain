package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type SubAccountTransferRequest struct {
	FromEmail       string      `json:"fromEmail"`
	ToEmail         string      `json:"toEmail"`
	Asset           oc.SymbolId `json:"asset"`
	Amount          oc.Amount   `json:"amount"`
	TimestampMillis int64       `json:"timestamp,omitempty"`
}

type SubAccountTransferResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	TxnId   string `json:"txnId"`
}

// SubAccountTransfer executes an asset transfer between the master account and a sub-account
// https://docs.binance.us/#execute-sub-account-transfer
func (c *Client) SubAccountTransfer(args *SubAccountTransferRequest) (*SubAccountTransferResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args cannot be nil")
	}

	var response SubAccountTransferResponse
	query := url.Values{}

	// Required parameters
	query.Set("fromEmail", args.FromEmail)
	query.Set("toEmail", args.ToEmail)
	query.Set("asset", string(args.Asset))
	query.Set("amount", args.Amount.String())
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	_, err := c.Request("POST", "/sapi/v3/sub-account/transfer", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
