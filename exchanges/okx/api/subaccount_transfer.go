package api

import oc "github.com/cordialsys/offchain"

type SubaccountTransferRequest struct {
	Currency       oc.SymbolId    `json:"ccy"`
	Amount         oc.Amount      `json:"amt"`
	From           oc.AccountType `json:"from"`
	To             oc.AccountType `json:"to"`
	FromSubAccount oc.AccountId   `json:"fromSubAccount"`
	ToSubAccount   oc.AccountId   `json:"toSubAccount"`
	LoanTrans      *bool          `json:"loanTrans,omitempty"`
	OmitPosRisk    *string        `json:"omitPosRisk,omitempty"`
}

type SubaccountTransferResponse Response[[]SubaccountTransferResult]
type SubaccountTransferResult struct {
	TransferID string `json:"transId,omitempty"`
}

// Rate limit: 1 request per second
// https://www.okx.com/docs-v5/en/#sub-account-rest-api-get-history-of-managed-sub-account-transfer
func (c *Client) SubaccountTransfer(args *SubaccountTransferRequest) (*SubaccountTransferResponse, error) {
	var response SubaccountTransferResponse
	_, err := c.Request("POST", "/api/v5/asset/subaccount/transfer", args, &response, nil)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
