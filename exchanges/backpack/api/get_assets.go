package api

import (
	oc "github.com/cordialsys/offchain"
)

type Token struct {
	Blockchain      oc.NetworkId       `json:"blockchain"`
	ContractAddress oc.ContractAddress `json:"contractAddress"`
	DepositEnabled  bool               `json:"depositEnabled"`
	WithdrawEnabled bool               `json:"withdrawEnabled"`

	MinimumDeposit oc.Amount `json:"minimumDeposit"`

	MinimumWithdraw   oc.Amount `json:"minimumWithdraw"`
	MaximumWithdrawal oc.Amount `json:"maximumWithdrawal"`

	WithdrawalFee oc.Amount `json:"withdrawalFee"`
}

type Asset struct {
	Symbol string   `json:"symbol"`
	Tokens []*Token `json:"tokens"`
}

// https://docs.backpack.exchange/#tag/Assets/operation/get_assets
func (c *Client) GetAssets() ([]Asset, error) {
	var response []Asset

	// No instruction type is specified in the docs, using a generic "assetQuery"
	// This might need to be updated based on actual API requirements
	_, err := c.Request("GET", "/api/v1/assets", "assetQuery", nil, &response, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
