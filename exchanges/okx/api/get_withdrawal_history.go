package api

import (
	"fmt"
	"net/url"

	oc "github.com/cordialsys/offchain"
	"github.com/cordialsys/offchain/client"
)

type WithdrawalHistoryRequest struct {
	Currency *oc.SymbolId `json:"ccy,omitempty"`
	WdId     string       `json:"wdId,omitempty"`
	ClientId string       `json:"clientId,omitempty"`
	TxId     string       `json:"txId,omitempty"`
	Type     string       `json:"type,omitempty"`
	State    string       `json:"state,omitempty"`
	After    int64        `json:"after,omitempty"`
	Before   int64        `json:"before,omitempty"`
	Limit    int          `json:"limit,omitempty"`
}

type AddressExtras map[string]string

type WithdrawalRecord struct {
	Currency         oc.SymbolId          `json:"ccy"`
	Chain            SymbolAndChain       `json:"chain"`
	NonTradableAsset bool                 `json:"nonTradableAsset"`
	Amount           oc.Amount            `json:"amt"`
	Timestamp        string               `json:"ts"`
	From             string               `json:"from"`
	AreaCodeFrom     string               `json:"areaCodeFrom"`
	To               string               `json:"to"`
	AreaCodeTo       string               `json:"areaCodeTo"`
	Tag              string               `json:"tag,omitempty"`
	PaymentId        string               `json:"pmtId,omitempty"`
	Memo             string               `json:"memo,omitempty"`
	AddressExtras    AddressExtras        `json:"addrEx,omitempty"`
	TxId             client.TransactionId `json:"txId"`
	Fee              oc.Amount            `json:"fee"`
	FeeCurrency      oc.SymbolId          `json:"feeCcy"`
	State            WithdrawalState      `json:"state"`
	WithdrawalId     string               `json:"wdId"`
	ClientId         string               `json:"clientId"`
	Note             string               `json:"note"`
}

type WithdrawalState string

const (
	WithdrawalStateCanceled          WithdrawalState = "-2"
	WithdrawalStateFailed            WithdrawalState = "-1"
	WithdrawalStateWaitingWithdrawal WithdrawalState = "0"
	WithdrawalStateBroadcasting      WithdrawalState = "1"
	WithdrawalStateSuccess           WithdrawalState = "2"
)

type WithdrawalHistoryResponse = Response[[]WithdrawalRecord]

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-cancel-withdrawal
func (c *Client) GetWithdrawalHistory(args *WithdrawalHistoryRequest) (*WithdrawalHistoryResponse, error) {
	var response WithdrawalHistoryResponse
	query := url.Values{}

	if args != nil {
		if args.Currency != nil {
			query.Set("ccy", string(*args.Currency))
		}
		if args.WdId != "" {
			query.Set("wdId", args.WdId)
		}
		if args.ClientId != "" {
			query.Set("clientId", args.ClientId)
		}
		if args.TxId != "" {
			query.Set("txId", args.TxId)
		}
		if args.Type != "" {
			query.Set("type", args.Type)
		}
		if args.State != "" {
			query.Set("state", args.State)
		}
		if args.After != 0 {
			query.Set("after", fmt.Sprintf("%d", args.After))
		}
		if args.Before != 0 {
			query.Set("before", fmt.Sprintf("%d", args.Before))
		}
		if args.Limit != 0 {
			query.Set("limit", fmt.Sprintf("%d", args.Limit))
		}
	}

	_, err := c.Request("GET", "/api/v5/asset/withdrawal-history", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
