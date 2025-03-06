package api

import (
	"fmt"

	rebal "github.com/cordialsys/offchain"
)

type BalanceResponse = Response[BalanceInfo]

type BalanceInfo struct {
	AvailableBalance rebal.Amount   `json:"availBal"`
	Balance          rebal.Amount   `json:"bal"`
	Currency         rebal.SymbolId `json:"ccy"`
	FrozenBalance    rebal.Amount   `json:"frozenBal"`
}

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-get-balance
func (c *Client) GetBalance(account rebal.GetBalanceArgs) (*BalanceResponse, error) {
	var response BalanceResponse
	_, err := c.Request("GET", "/api/v5/asset/balances", nil, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &response, nil
}
