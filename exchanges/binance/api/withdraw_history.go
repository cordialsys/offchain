package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	oc "github.com/cordialsys/offchain"
)

// WithdrawalStatus represents the status of a withdrawal
type WithdrawalStatus int

const (
	WithdrawalStatusEmailSent        WithdrawalStatus = 0
	WithdrawalStatusAwaitingApproval WithdrawalStatus = 2
	WithdrawalStatusRejected         WithdrawalStatus = 3
	WithdrawalStatusProcessing       WithdrawalStatus = 4
	WithdrawalStatusCompleted        WithdrawalStatus = 6
)

// String returns the string representation of the withdrawal status
func (s WithdrawalStatus) String() string {
	switch s {
	case WithdrawalStatusEmailSent:
		return "Email Sent"
	case WithdrawalStatusAwaitingApproval:
		return "Awaiting Approval"
	case WithdrawalStatusRejected:
		return "Rejected"
	case WithdrawalStatusProcessing:
		return "Processing"
	case WithdrawalStatusCompleted:
		return "Completed"
	default:
		return "Unknown"
	}
}

const (
	// MaxWithdrawalHistoryLimit     = 1000
	DefaultWithdrawalHistoryLimit = 1000
	MaxIdListLength               = 45
	MaxTimeRangeInDays            = 90
	MaxTimeRangeWithOrderIdDays   = 7
)

type TransferType int

const (
	TransferTypeExternal TransferType = 0
	TransferTypeInternal TransferType = 1
)

type WithdrawalHistoryRequest struct {
	Coin            *oc.SymbolId      `json:"coin,omitempty"`
	WithdrawOrderId string            `json:"withdrawOrderId,omitempty"`
	Status          *WithdrawalStatus `json:"status,omitempty"`
	Offset          *int              `json:"offset,omitempty"`
	Limit           *int              `json:"limit,omitempty"`
	IdList          []string          `json:"idList,omitempty"`
	StartTime       *int64            `json:"startTime,omitempty"`
	EndTime         *int64            `json:"endTime,omitempty"`
}

type WithdrawalRecord struct {
	Id              string           `json:"id"`
	Amount          oc.Amount        `json:"amount"`
	TransactionFee  oc.Amount        `json:"transactionFee"`
	Coin            oc.SymbolId      `json:"coin"`
	Status          WithdrawalStatus `json:"status"`
	Address         oc.Address       `json:"address"`
	TxId            oc.TransactionId `json:"txId"`
	ApplyTime       string           `json:"applyTime"`
	Network         oc.NetworkId     `json:"network"`
	TransferType    TransferType     `json:"transferType"`
	WithdrawOrderId string           `json:"withdrawOrderId,omitempty"`
	Info            string           `json:"info"`
	ConfirmNo       int              `json:"confirmNo"`
	WalletType      WalletType       `json:"walletType"`
	TxKey           string           `json:"txKey"`
	CompleteTime    string           `json:"completeTime,omitempty"`
}

// https://developers.binance.com/docs/wallet/capital/withdraw-history
func (c *Client) GetWithdrawalHistory(args *WithdrawalHistoryRequest) ([]WithdrawalRecord, error) {
	var response []WithdrawalRecord
	query := url.Values{}
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

	if args != nil {
		// Validate idList length
		if len(args.IdList) > MaxIdListLength {
			return nil, fmt.Errorf("idList cannot contain more than %d IDs", MaxIdListLength)
		}

		// Validate and set limit
		if args.Limit != nil {
			query.Set("limit", strconv.Itoa(*args.Limit))
		}

		// Set other parameters
		if args.Coin != nil {
			query.Set("coin", string(*args.Coin))
		}
		if args.WithdrawOrderId != "" {
			query.Set("withdrawOrderId", args.WithdrawOrderId)
		}
		if args.Status != nil {
			query.Set("status", strconv.Itoa(int(*args.Status)))
		}
		if args.Offset != nil {
			query.Set("offset", strconv.Itoa(*args.Offset))
		}
		if len(args.IdList) > 0 {
			query.Set("idList", strings.Join(args.IdList, ","))
		}
		if args.StartTime != nil {
			query.Set("startTime", strconv.FormatInt(*args.StartTime, 10))
		}
		if args.EndTime != nil {
			query.Set("endTime", strconv.FormatInt(*args.EndTime, 10))
		}
	}

	_, err := c.Request("GET", "/sapi/v1/capital/withdraw/history", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return response, nil
}
