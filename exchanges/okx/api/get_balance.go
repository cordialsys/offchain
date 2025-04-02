package api

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
)

type BalanceResponse = Response[[]BalanceInfo]

type BalanceInfo struct {
	AvailableBalance oc.Amount   `json:"availBal"`
	Balance          oc.Amount   `json:"bal"`
	Currency         oc.SymbolId `json:"ccy"`
	FrozenBalance    oc.Amount   `json:"frozenBal"`
}

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-get-balance
func (c *Client) GetBalance(account client.GetBalanceArgs) (*BalanceResponse, error) {
	var response BalanceResponse
	_, err := c.Request("GET", "/api/v5/asset/balances", nil, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}
	return &response, nil
}
