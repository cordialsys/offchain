package api

import (
	oc "github.com/cordialsys/offchain"
)

type NetworkInfo struct {
	AddressRegex         string             `json:"addressRegex"`
	Coin                 string             `json:"coin"`
	DepositDesc          string             `json:"depositDesc,omitempty"`
	DepositEnable        bool               `json:"depositEnable"`
	IsDefault            bool               `json:"isDefault"`
	MemoRegex            string             `json:"memoRegex"`
	MinConfirm           int                `json:"minConfirm"`
	Name                 string             `json:"name"`
	Network              oc.NetworkId       `json:"network"`
	SpecialTips          string             `json:"specialTips"`
	UnLockConfirm        int                `json:"unLockConfirm"`
	WithdrawDesc         string             `json:"withdrawDesc,omitempty"`
	WithdrawEnable       bool               `json:"withdrawEnable"`
	WithdrawFee          string             `json:"withdrawFee"`
	WithdrawMin          string             `json:"withdrawMin"`
	WithdrawMax          string             `json:"withdrawMax"`
	WithdrawInternalMin  string             `json:"withdrawInternalMin"`
	SameAddress          bool               `json:"sameAddress"`
	EstimatedArrivalTime int                `json:"estimatedArrivalTime"`
	Busy                 bool               `json:"busy"`
	ContractAddressUrl   string             `json:"contractAddressUrl"`
	ContractAddress      oc.ContractAddress `json:"contractAddress"`
	Denomination         int64              `json:"denomination,omitempty"`
}

type CoinInformation struct {
	Coin              oc.SymbolId   `json:"coin"`
	DepositAllEnable  bool          `json:"depositAllEnable"`
	Free              oc.Amount     `json:"free"`
	Freeze            oc.Amount     `json:"freeze"`
	Ipoable           oc.Amount     `json:"ipoable"`
	Ipoing            oc.Amount     `json:"ipoing"`
	IsLegalMoney      bool          `json:"isLegalMoney"`
	Locked            oc.Amount     `json:"locked"`
	Name              string        `json:"name"`
	NetworkList       []NetworkInfo `json:"networkList"`
	Storage           oc.Amount     `json:"storage"`
	Trading           bool          `json:"trading"`
	WithdrawAllEnable bool          `json:"withdrawAllEnable"`
	Withdrawing       oc.Amount     `json:"withdrawing"`
}

// https://developers.binance.com/docs/wallet/capital
func (c *Client) GetAllCoinsInformation() ([]CoinInformation, error) {
	var response []CoinInformation
	_, err := c.Request("GET", "/sapi/v1/capital/config/getall", nil, &response, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}
