package api

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
)

// WithdrawalRequest represents the request parameters for a withdrawal
type WithdrawalRequest struct {
	Address        oc.Address   `json:"address"`
	Blockchain     oc.NetworkId `json:"blockchain"`
	ClientId       string       `json:"clientId,omitempty"`
	Quantity       oc.Amount    `json:"quantity"`
	Symbol         oc.SymbolId  `json:"symbol"`
	TwoFactorToken string       `json:"twoFactorToken,omitempty"`
	AutoBorrow     *bool        `json:"autoBorrow,omitempty"`
	AutoLendRedeem *bool        `json:"autoLendRedeem,omitempty"`
}

// WithdrawalResponse represents the response from the withdrawal endpoint
type WithdrawalResponse struct {
	ID              int                  `json:"id"`
	Blockchain      string               `json:"blockchain"`
	ClientId        string               `json:"clientId,omitempty"`
	Identifier      string               `json:"identifier,omitempty"`
	Quantity        oc.Amount            `json:"quantity"`
	Fee             oc.Amount            `json:"fee"`
	Symbol          oc.SymbolId          `json:"symbol"`
	Status          string               `json:"status"`
	SubaccountId    int                  `json:"subaccountId,omitempty"`
	ToAddress       oc.Address           `json:"toAddress"`
	TransactionHash client.TransactionId `json:"transactionHash,omitempty"`
	CreatedAt       string               `json:"createdAt"`
	IsInternal      bool                 `json:"isInternal"`
}

// https://docs.backpack.exchange/#tag/Capital/operation/request_withdrawal
func (c *Client) RequestWithdrawal(req *WithdrawalRequest) (*WithdrawalResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("withdrawal request cannot be nil")
	}

	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}

	if req.Blockchain == "" {
		return nil, fmt.Errorf("blockchain is required")
	}

	if req.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	var response WithdrawalResponse
	_, err := c.Request("POST", "/wapi/v1/capital/withdrawals", "withdraw", req, &response, nil)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
