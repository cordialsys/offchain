package api

import (
	"fmt"
	"net/url"
	"time"

	oc "github.com/cordialsys/offchain"
)

type WithdrawRequest struct {
	Coin               oc.SymbolId  `json:"coin"`
	WithdrawOrderId    string       `json:"withdrawOrderId,omitempty"`
	Network            oc.NetworkId `json:"network,omitempty"`
	Address            oc.Address   `json:"address"`
	AddressTag         string       `json:"addressTag,omitempty"`
	Amount             oc.Amount    `json:"amount"`
	TransactionFeeFlag *bool        `json:"transactionFeeFlag,omitempty"`
	Name               string       `json:"name,omitempty"`
	// The wallet type for withdraw，0-spot wallet ，1-funding wallet. Default walletType is the current "selected wallet" under wallet->Fiat and Spot/Funding->Deposit
	WalletType *WalletType `json:"walletType,omitempty"`
}

type WithdrawResponse struct {
	Id string `json:"id"`
}

type WalletType int

const WalletTypeSpot WalletType = 0
const WalletTypeFunding WalletType = 1

// https://developers.binance.com/docs/wallet/capital/withdraw
func (c *Client) Withdraw(args *WithdrawRequest) (*WithdrawResponse, error) {
	var response WithdrawResponse
	query := url.Values{}

	// Required parameters
	query.Set("coin", string(args.Coin))
	query.Set("address", string(args.Address))
	query.Set("amount", args.Amount.String())
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

	// Optional parameters
	if args.WithdrawOrderId != "" {
		query.Set("withdrawOrderId", args.WithdrawOrderId)
	}
	if args.Network != "" {
		query.Set("network", string(args.Network))
	}
	if args.AddressTag != "" {
		query.Set("addressTag", args.AddressTag)
	}
	if args.TransactionFeeFlag != nil {
		if *args.TransactionFeeFlag {
			query.Set("transactionFeeFlag", "true")
		} else {
			query.Set("transactionFeeFlag", "false")
		}
	}
	if args.Name != "" {
		query.Set("name", args.Name)
	}
	if args.WalletType != nil {
		query.Set("walletType", fmt.Sprintf("%d", *args.WalletType))
	}

	_, err := c.Request("POST", "/sapi/v1/capital/withdraw/apply", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
