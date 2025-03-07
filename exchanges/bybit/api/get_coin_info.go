package api

import (
	"net/url"

	oc "github.com/cordialsys/offchain"
)

type GetCoinInfoResponse = Response[GetCoinInfoResult]

type GetCoinInfoResult struct {
	Rows []CoinInfo `json:"rows"`
}

type CoinInfo struct {
	Name         string  `json:"name"`
	Coin         string  `json:"coin"`
	RemainAmount string  `json:"remainAmount"`
	Chains       []Chain `json:"chains"`
}

type Chain struct {
	ChainType             string `json:"chainType"`
	Confirmation          string `json:"confirmation"`
	WithdrawFee           string `json:"withdrawFee"`
	DepositMin            string `json:"depositMin"`
	WithdrawMin           string `json:"withdrawMin"`
	Chain                 string `json:"chain"`
	ChainDeposit          string `json:"chainDeposit"`
	ChainWithdraw         string `json:"chainWithdraw"`
	MinAccuracy           string `json:"minAccuracy"`
	WithdrawPercentageFee string `json:"withdrawPercentageFee"`
	ContractAddress       string `json:"contractAddress"`
}

// https://bybit-exchange.github.io/docs/v5/asset/coin-info/
func (c *Client) GetCoinInfo(coinMaybe oc.SymbolId) (*GetCoinInfoResponse, error) {
	params := url.Values{}
	if coinMaybe != "" {
		params.Set("coin", string(coinMaybe))
	}
	var response GetCoinInfoResponse
	_, err := c.Request("GET", "/v5/asset/coin/query-info", nil, &response, params)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
