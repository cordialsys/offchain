package api

import (
	"fmt"
	"net/url"
	"time"

	oc "github.com/cordialsys/offchain"
)

type AccountType string

const (
	AccountTypeSpot           AccountType = "SPOT"
	AccountTypeUsdtFuture     AccountType = "USDT_FUTURE"
	AccountTypeCoinFuture     AccountType = "COIN_FUTURE"
	AccountTypeMargin         AccountType = "MARGIN"
	AccountTypeIsolatedMargin AccountType = "ISOLATED_MARGIN"
)

type UniversalTransferRequest struct {
	FromEmail       oc.AccountId `json:"fromEmail,omitempty"`
	ToEmail         oc.AccountId `json:"toEmail,omitempty"`
	FromAccountType AccountType  `json:"fromAccountType"`
	ToAccountType   AccountType  `json:"toAccountType"`
	ClientTranId    string       `json:"clientTranId,omitempty"`
	Symbol          oc.SymbolId  `json:"symbol,omitempty"`
	Asset           oc.SymbolId  `json:"asset"`
	Amount          oc.Amount    `json:"amount"`
}

type UniversalTransferResponse struct {
	TranId       int64  `json:"tranId"`
	ClientTranId string `json:"clientTranId,omitempty"`
}

// UniversalTransfer performs a universal transfer between accounts
// https://developers.binance.com/docs/sub_account/asset-management/Universal-Transfer
func (c *Client) UniversalTransfer(args *UniversalTransferRequest) (*UniversalTransferResponse, error) {
	var response UniversalTransferResponse
	query := url.Values{}

	// Required parameters
	query.Set("fromAccountType", string(args.FromAccountType))
	query.Set("toAccountType", string(args.ToAccountType))
	query.Set("asset", string(args.Asset))
	query.Set("amount", args.Amount.String())
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

	// Optional parameters
	if args.FromEmail != "" {
		query.Set("fromEmail", string(args.FromEmail))
	}
	if args.ToEmail != "" {
		query.Set("toEmail", string(args.ToEmail))
	}
	if args.ClientTranId != "" {
		query.Set("clientTranId", args.ClientTranId)
	}
	if args.Symbol != "" {
		query.Set("symbol", string(args.Symbol))
	}

	_, err := c.Request("POST", "/sapi/v1/sub-account/universalTransfer", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
