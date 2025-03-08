package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type WithdrawalStatus int

const (
	WithdrawalStatusEmailSent        WithdrawalStatus = 0
	WithdrawalStatusCanceled         WithdrawalStatus = 1
	WithdrawalStatusAwaitingApproval WithdrawalStatus = 2
	WithdrawalStatusRejected         WithdrawalStatus = 3
	WithdrawalStatusProcessing       WithdrawalStatus = 4
	WithdrawalStatusFailure          WithdrawalStatus = 5
	WithdrawalStatusCompleted        WithdrawalStatus = 6
)

func (s WithdrawalStatus) String() string {
	switch s {
	case WithdrawalStatusEmailSent:
		return "EMAIL_SENT"
	case WithdrawalStatusCanceled:
		return "CANCELED"
	case WithdrawalStatusAwaitingApproval:
		return "AWAITING_APPROVAL"
	case WithdrawalStatusRejected:
		return "REJECTED"
	case WithdrawalStatusProcessing:
		return "PROCESSING"
	case WithdrawalStatusFailure:
		return "FAILURE"
	case WithdrawalStatusCompleted:
		return "COMPLETED"
	default:
		return "UNKNOWN"
	}
}

type GetCryptoWithdrawalHistoryRequest struct {
	Coin            *oc.SymbolId      `json:"coin,omitempty"`
	WithdrawOrderId string            `json:"withdrawOrderId,omitempty"`
	Status          *WithdrawalStatus `json:"status,omitempty"`
	StartTime       *int64            `json:"startTime,omitempty"`
	EndTime         *int64            `json:"endTime,omitempty"`
	Offset          *int              `json:"offset,omitempty"`
	Limit           *int              `json:"limit,omitempty"`
	RecvWindow      *int64            `json:"recvWindow,omitempty"`
	TimestampMillis int64             `json:"timestamp,omitempty"`
}

type CryptoWithdrawalRecord struct {
	Id             string           `json:"id"`
	Amount         oc.Amount        `json:"amount"`
	TransactionFee oc.Amount        `json:"transactionFee"`
	Coin           oc.SymbolId      `json:"coin"`
	Status         WithdrawalStatus `json:"status"`
	Address        oc.Address       `json:"address"`
	ApplyTime      string           `json:"applyTime"`
	Network        oc.NetworkId     `json:"network"`
	TransferType   int              `json:"transferType"`
	TxId           oc.TransactionId `json:"txId"`
}

type FiatWithdrawalHistoryRequest struct {
	FiatCurrency   string `json:"fiatCurrency,omitempty"`
	OrderId        string `json:"orderId,omitempty"`
	Offset         *int   `json:"offset,omitempty"`
	PaymentChannel string `json:"paymentChannel,omitempty"`
	PaymentMethod  string `json:"paymentMethod,omitempty"`
	StartTime      *int64 `json:"startTime,omitempty"`
	EndTime        *int64 `json:"endTime,omitempty"`
	RecvWindow     *int64 `json:"recvWindow,omitempty"`
}

type FiatWithdrawalRecord struct {
	OrderId        string `json:"orderId"`
	PaymentAccount string `json:"paymentAccount"`
	PaymentChannel string `json:"paymentChannel"`
	PaymentMethod  string `json:"paymentMethod"`
	OrderStatus    string `json:"orderStatus"`
	Amount         string `json:"amount"`
	TransactionFee string `json:"transactionFee"`
	PlatformFee    string `json:"platformFee"`
}

type FiatWithdrawalHistoryResponse struct {
	AssetLogRecordList []FiatWithdrawalRecord `json:"assetLogRecordList"`
}

// GetCryptoWithdrawalHistory retrieves crypto withdrawal history
func (c *Client) GetCryptoWithdrawalHistory(args *GetCryptoWithdrawalHistoryRequest) ([]CryptoWithdrawalRecord, error) {
	var response []CryptoWithdrawalRecord
	query := url.Values{}
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))
	if args != nil {
		if args.Coin != nil {
			query.Set("coin", string(*args.Coin))
		}
		if args.WithdrawOrderId != "" {
			query.Set("withdrawOrderId", args.WithdrawOrderId)
		}
		if args.Status != nil {
			query.Set("status", strconv.Itoa(int(*args.Status)))
		}
		if args.StartTime != nil {
			query.Set("startTime", strconv.FormatInt(*args.StartTime, 10))
		}
		if args.EndTime != nil {
			query.Set("endTime", strconv.FormatInt(*args.EndTime, 10))
		}
		if args.Offset != nil {
			query.Set("offset", strconv.Itoa(*args.Offset))
		}
		if args.Limit != nil {
			if *args.Limit > 1000 {
				return nil, fmt.Errorf("limit cannot exceed 1000")
			}
			query.Set("limit", strconv.Itoa(*args.Limit))
		}
		if args.RecvWindow != nil {
			if *args.RecvWindow > 60000 {
				return nil, fmt.Errorf("recvWindow cannot be greater than 60000")
			}
			query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
		}
	}

	_, err := c.Request("GET", "/sapi/v1/capital/withdraw/history", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetFiatWithdrawalHistory retrieves fiat withdrawal history
func (c *Client) GetFiatWithdrawalHistory(args *FiatWithdrawalHistoryRequest) (*FiatWithdrawalHistoryResponse, error) {
	var response FiatWithdrawalHistoryResponse
	query := url.Values{}

	if args != nil {
		if args.FiatCurrency != "" {
			query.Set("fiatCurrency", args.FiatCurrency)
		}
		if args.OrderId != "" {
			query.Set("orderId", args.OrderId)
		}
		if args.Offset != nil {
			query.Set("offset", strconv.Itoa(*args.Offset))
		}
		if args.PaymentChannel != "" {
			query.Set("paymentChannel", args.PaymentChannel)
		}
		if args.PaymentMethod != "" {
			query.Set("paymentMethod", args.PaymentMethod)
		}
		if args.StartTime != nil {
			query.Set("startTime", strconv.FormatInt(*args.StartTime, 10))
		}
		if args.EndTime != nil {
			query.Set("endTime", strconv.FormatInt(*args.EndTime, 10))
		}
		if args.RecvWindow != nil {
			if *args.RecvWindow > 60000 {
				return nil, fmt.Errorf("recvWindow cannot be greater than 60000")
			}
			query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
		}
	}

	_, err := c.Request("GET", "/sapi/v1/fiatpayment/query/withdraw/history", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
