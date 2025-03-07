package api

import (
	oc "github.com/cordialsys/offchain"
)

type WithdrawRequest struct {
	Coin        oc.SymbolId  `json:"coin"`                  // Required, uppercase only
	Chain       oc.NetworkId `json:"chain,omitempty"`       // Optional, depends on forceChain
	Address     oc.Address   `json:"address"`               // Required, wallet address or Bybit UID
	Tag         string       `json:"tag,omitempty"`         // Optional, required if tag exists in address book
	Amount      oc.Amount    `json:"amount"`                // Required, withdrawal amount
	Timestamp   int64        `json:"timestamp"`             // Required, current timestamp in ms
	ForceChain  *int         `json:"forceChain,omitempty"`  // Optional, 0=default, 1=force on-chain, 2=UID withdraw
	AccountType AccountType  `json:"accountType,omitempty"` // Optional, SPOT(default) or FUND
	FeeType     *int         `json:"feeType,omitempty"`     // Optional, 0=manual fee calc, 1=auto deduct
	RequestID   string       `json:"requestId,omitempty"`   // Optional, unique ID for idempotency
	Beneficiary *Beneficiary `json:"beneficiary,omitempty"` // Optional, travel rule info
}
type Beneficiary struct {
	// Add travel rule fields here as needed
	// This would depend on the specific requirements for KOR/TR/KZ users
}

type WithdrawResult struct {
	ID string `json:"id"` // Withdrawal ID
}
type WithdrawResponse = Response[WithdrawResult]

// https://bybit-exchange.github.io/docs/v5/asset/withdraw
func (c *Client) Withdraw(request *WithdrawRequest) (*WithdrawResponse, error) {
	var response WithdrawResponse
	_, err := c.Request("POST", "/v5/asset/withdraw/create", request, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
