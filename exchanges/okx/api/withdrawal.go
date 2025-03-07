package api

import oc "github.com/cordialsys/offchain"

type WithdrawalRequest struct {
	Amount         string      `json:"amt"`
	Destination    string      `json:"dest"`
	Currency       oc.SymbolId `json:"ccy"`
	SymbolAndChain string      `json:"chain"`
	ToAddress      oc.Address  `json:"toAddr"`
}

type WithdrawalResponse Response[[]WithdrawalResult]

type WithdrawalResult struct {
	WithdrawalID string `json:"wdId,omitempty"`
	Currency     string `json:"ccy,omitempty"`
	Chain        string `json:"chain,omitempty"`
	Amount       string `json:"amt,omitempty"`
	ToAddress    string `json:"toAddr,omitempty"`
	TxID         string `json:"txId,omitempty"`
	State        string `json:"state,omitempty"`
}

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-withdrawal
func (c *Client) Withdrawal(args *WithdrawalRequest) (*WithdrawalResponse, error) {
	var response WithdrawalResponse
	_, err := c.Request("POST", "/api/v5/asset/withdrawal", args, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
