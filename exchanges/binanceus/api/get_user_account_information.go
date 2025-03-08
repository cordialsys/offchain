package api

import (
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type CommissionRates struct {
	Maker  string `json:"maker"`
	Taker  string `json:"taker"`
	Buyer  string `json:"buyer"`
	Seller string `json:"seller"`
}

type Balance struct {
	Asset  oc.SymbolId `json:"asset"`
	Free   oc.Amount   `json:"free"`
	Locked oc.Amount   `json:"locked"`
}

type AccountType string

const (
	AccountTypeSpot AccountType = "SPOT"
)

type AccountResponse struct {
	MakerCommission            int             `json:"makerCommission"`
	TakerCommission            int             `json:"takerCommission"`
	BuyerCommission            int             `json:"buyerCommission"`
	SellerCommission           int             `json:"sellerCommission"`
	CommissionRates            CommissionRates `json:"commissionRates"`
	CanTrade                   bool            `json:"canTrade"`
	CanWithdraw                bool            `json:"canWithdraw"`
	CanDeposit                 bool            `json:"canDeposit"`
	Brokered                   bool            `json:"brokered"`
	RequireSelfTradePrevention bool            `json:"requireSelfTradePrevention"`
	UpdateTime                 int64           `json:"updateTime"`
	AccountType                AccountType     `json:"accountType"`
	Balances                   []Balance       `json:"balances"`
	Permissions                []string        `json:"permissions"`
}

type GetAccountRequest struct {
	RecvWindow      *int64 `json:"recvWindow,omitempty"`
	TimestampMillis int64  `json:"timestamp,omitempty"`
}

// https://docs.binance.us/#get-current-account-information-user_data
func (c *Client) GetAccount(args *GetAccountRequest) (*AccountResponse, error) {
	var response AccountResponse
	query := url.Values{}

	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	if args != nil && args.RecvWindow != nil {
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("GET", "/api/v3/account", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
