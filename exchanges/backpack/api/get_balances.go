package api

import (
	oc "github.com/cordialsys/offchain"
)

// Balance represents the balance information for a specific asset
type Balance struct {
	Available oc.Amount `json:"available"`
	Locked    oc.Amount `json:"locked"`
	Staked    oc.Amount `json:"staked"`
}

type BalanceResponse map[string]Balance

// https://docs.backpack.exchange/#tag/Capital/operation/get_balances
func (c *Client) GetBalances() (BalanceResponse, error) {
	var response BalanceResponse

	// No query parameters or request body needed for this endpoint
	_, err := c.Request("GET", "/api/v1/capital", "balanceQuery", nil, &response, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
