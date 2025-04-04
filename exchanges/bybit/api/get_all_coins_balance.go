package api

import (
	"net/url"

	oc "github.com/cordialsys/offchain"
)

type GetAllCoinsBalanceResponse = Response[GetAllCoinsBalanceResult]

type GetAllCoinsBalanceResult struct {
	MemberID    string    `json:"memberId"`
	AccountType string    `json:"accountType"`
	Balance     []Balance `json:"balance"`
}

type Balance struct {
	Coin            oc.SymbolId `json:"coin"`
	TransferBalance oc.Amount   `json:"transferBalance"`
	WalletBalance   oc.Amount   `json:"walletBalance"`
	Bonus           string      `json:"bonus"`
}

type AccountType = string

const (
	AccountTypeUnified AccountType = "UNIFIED"
	AccountTypeFund    AccountType = "FUND"
	AccountTypeSpot    AccountType = "SPOT"
	AccountTypeLinear  AccountType = "LINEAR"
	AccountTypeInverse AccountType = "INVERSE"
)

// https://bybit-exchange.github.io/docs/v5/asset/balance/all-balance
func (c *Client) GetAllCoinsBalance(accountType oc.AccountType, coinMaybe oc.SymbolId) (*GetAllCoinsBalanceResponse, error) {

	params := url.Values{}
	params.Set("accountType", string(accountType))
	if coinMaybe != "" {
		params.Set("coin", string(coinMaybe))
	}
	var response GetAllCoinsBalanceResponse
	_, err := c.Request("GET", "/v5/asset/transfer/query-account-coins-balance", nil, &response, params)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
