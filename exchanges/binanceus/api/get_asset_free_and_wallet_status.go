package api

import (
	"fmt"
	"net/url"
	"strconv"

	oc "github.com/cordialsys/offchain"
)

type NetworkConfig struct {
	Network                 oc.NetworkId `json:"network"`
	Coin                    oc.SymbolId  `json:"coin"`
	WithdrawIntegerMultiple string       `json:"withdrawIntegerMultiple"`
	IsDefault               bool         `json:"isDefault"`
	DepositEnable           bool         `json:"depositEnable"`
	WithdrawEnable          bool         `json:"withdrawEnable"`
	DepositDesc             string       `json:"depositDesc"`
	WithdrawDesc            string       `json:"withdrawDesc"`
	Name                    string       `json:"name"`
	ResetAddressStatus      bool         `json:"resetAddressStatus"`
	AddressRegex            string       `json:"addressRegex,omitempty"`
	MemoRegex               string       `json:"memoRegex,omitempty"`
	WithdrawFee             string       `json:"withdrawFee"`
	WithdrawMin             string       `json:"withdrawMin"`
	WithdrawMax             string       `json:"withdrawMax"`
	MinConfirm              int          `json:"minConfirm,omitempty"`
	UnLockConfirm           int          `json:"unLockConfirm,omitempty"`
	SpecialTips             string       `json:"specialTips,omitempty"`
}

type AssetConfig struct {
	Coin              oc.SymbolId     `json:"coin"`
	DepositAllEnable  bool            `json:"depositAllEnable"`
	WithdrawAllEnable bool            `json:"withdrawAllEnable"`
	Name              string          `json:"name"`
	Free              oc.Amount       `json:"free"`
	Locked            oc.Amount       `json:"locked"`
	Freeze            oc.Amount       `json:"freeze"`
	Withdrawing       oc.Amount       `json:"withdrawing"`
	Ipoing            oc.Amount       `json:"ipoing"`
	Ipoable           oc.Amount       `json:"ipoable"`
	Storage           oc.Amount       `json:"storage"`
	IsLegalMoney      bool            `json:"isLegalMoney"`
	Trading           bool            `json:"trading"`
	NetworkList       []NetworkConfig `json:"networkList"`
}

type GetAssetConfigRequest struct {
	RecvWindow      *int64 `json:"recvWindow,omitempty"`
	TimestampMillis int64  `json:"timestamp,omitempty"`
}

// https://docs.binance.us/#asset-fees-amp-wallet-status
func (c *Client) GetAssetConfig(args *GetAssetConfigRequest) ([]AssetConfig, error) {
	var response []AssetConfig
	query := url.Values{}
	query.Set("timestamp", strconv.FormatInt(args.TimestampMillis, 10))

	if args != nil && args.RecvWindow != nil {
		if *args.RecvWindow > 60000 {
			return nil, fmt.Errorf("recvWindow cannot be greater than 60000")
		}
		query.Set("recvWindow", strconv.FormatInt(*args.RecvWindow, 10))
	}

	_, err := c.Request("GET", "/sapi/v1/capital/config/getall", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return response, nil
}
