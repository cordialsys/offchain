package api

import oc "github.com/cordialsys/offchain"

type AccountTransferRequest struct {
	Currency oc.SymbolId `json:"ccy"`
	Amount   oc.Amount   `json:"amt"`
	From     string      `json:"from"`
	To       string      `json:"to"`
}

type AccountTransferResponse Response[AccountTransferResult]
type AccountTransferResult struct {
	TransferID string `json:"transId,omitempty"`
	Currency   string `json:"ccy,omitempty"`
	Amount     string `json:"amt,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-funds-transfer
func (c *Client) FundsTransfer(args *AccountTransferRequest) (*AccountTransferResponse, error) {
	var response AccountTransferResponse
	_, err := c.Request("POST", "/api/v5/asset/transfer", args, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
