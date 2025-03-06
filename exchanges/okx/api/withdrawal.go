package api

type WithdrawalRequest struct {
	Amount      string `json:"amt"`
	Destination string `json:"dest"`
	Currency    string `json:"ccy"`
	Chain       string `json:"chain"`
	ToAddress   string `json:"toAddr"`
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
func (c *Client) Withdrawal(args *WithdrawalRequest) (*AccountTransferResponse, error) {
	var response AccountTransferResponse
	_, err := c.Request("POST", "/api/v5/asset/withdrawal", args, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
