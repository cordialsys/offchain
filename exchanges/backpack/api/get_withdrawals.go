package api

import (
	"net/url"
	"strconv"
)

type WithdrawalHistoryRequest struct {
	From   *int64  `json:"from,omitempty"`
	To     *int64  `json:"to,omitempty"`
	Limit  *uint64 `json:"limit,omitempty"`
	Offset *uint64 `json:"offset,omitempty"`
}

// https://docs.backpack.exchange/#tag/Capital/operation/get_withdrawals
func (c *Client) GetWithdrawals(req *WithdrawalHistoryRequest) ([]WithdrawalResponse, error) {
	query := url.Values{}

	if req != nil {
		if req.From != nil {
			query.Set("from", strconv.FormatInt(*req.From, 10))
		}
		if req.To != nil {
			query.Set("to", strconv.FormatInt(*req.To, 10))
		}
		if req.Limit != nil {
			query.Set("limit", strconv.FormatUint(*req.Limit, 10))
		}
		if req.Offset != nil {
			query.Set("offset", strconv.FormatUint(*req.Offset, 10))
		}
	}

	var response []WithdrawalResponse
	_, err := c.Request("GET", "/wapi/v1/capital/withdrawals", "withdrawalQueryAll", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return response, nil
}
