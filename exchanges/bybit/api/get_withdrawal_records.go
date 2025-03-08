package api

import (
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type WithdrawalRecordsRequest struct {
	WithdrawID   string           `json:"withdrawID,omitempty"`
	TxID         oc.TransactionId `json:"txID,omitempty"`
	Coin         *oc.SymbolId     `json:"coin,omitempty"`
	WithdrawType *int             `json:"withdrawType,omitempty"`
	StartTime    *int64           `json:"startTime,omitempty"`
	EndTime      *int64           `json:"endTime,omitempty"`
	Limit        *int             `json:"limit,omitempty"`
	Cursor       string           `json:"cursor,omitempty"`
}

type WithdrawalRecord struct {
	WithdrawId   string           `json:"withdrawId"`
	TxID         oc.TransactionId `json:"txID"`
	WithdrawType int              `json:"withdrawType"`
	Coin         oc.SymbolId      `json:"coin"`
	Chain        oc.NetworkId     `json:"chain"`
	Amount       oc.Amount        `json:"amount"`
	WithdrawFee  oc.Amount        `json:"withdrawFee"`
	Status       WithdrawalStatus `json:"status"`
	ToAddress    oc.Address       `json:"toAddress"`
	Tag          string           `json:"tag"`
	CreateTime   string           `json:"createTime"`
	UpdateTime   string           `json:"updateTime"`
}

// WithdrawalStatus represents the status of a withdrawal
type WithdrawalStatus string

const (
	WithdrawalStatusSecurityCheck           WithdrawalStatus = "SecurityCheck"
	WithdrawalStatusPending                 WithdrawalStatus = "Pending"
	WithdrawalStatusSuccess                 WithdrawalStatus = "success"
	WithdrawalStatusCancelByUser            WithdrawalStatus = "CancelByUser"
	WithdrawalStatusReject                  WithdrawalStatus = "Reject"
	WithdrawalStatusFail                    WithdrawalStatus = "Fail"
	WithdrawalStatusBlockchainConfirmed     WithdrawalStatus = "BlockchainConfirmed"
	WithdrawalStatusMoreInformationRequired WithdrawalStatus = "MoreInformationRequired"
	WithdrawalStatusUnknown                 WithdrawalStatus = "Unknown"
)

type WithdrawalRecordsResult struct {
	Rows           []WithdrawalRecord `json:"rows"`
	NextPageCursor string             `json:"nextPageCursor"`
}

type WithdrawalRecordsResponse = Response[WithdrawalRecordsResult]

// https://bybit-exchange.github.io/docs/v5/asset/withdraw/withdraw-record
func (c *Client) GetWithdrawalRecords(args *WithdrawalRecordsRequest) (*WithdrawalRecordsResponse, error) {
	var response WithdrawalRecordsResponse
	query := url.Values{}

	if args != nil {
		if args.WithdrawID != "" {
			query.Set("withdrawID", args.WithdrawID)
		}
		if args.TxID != "" {
			query.Set("txID", string(args.TxID))
		}
		if args.Coin != nil {
			query.Set("coin", string(*args.Coin))
		}
		if args.WithdrawType != nil {
			query.Set("withdrawType", strconv.Itoa(*args.WithdrawType))
		}
		if args.StartTime != nil {
			query.Set("startTime", strconv.FormatInt(*args.StartTime, 10))
		}
		if args.EndTime != nil {
			query.Set("endTime", strconv.FormatInt(*args.EndTime, 10))
		}
		if args.Limit != nil {
			query.Set("limit", strconv.Itoa(*args.Limit))
		}
		if args.Cursor != "" {
			query.Set("cursor", args.Cursor)
		}
	}

	_, err := c.Request("GET", "/v5/asset/withdraw/query-record", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
