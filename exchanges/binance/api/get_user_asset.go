package api

import (
	"fmt"
	"net/url"
	"time"

	oc "github.com/cordialsys/offchain"
)

type UserAssetRequest struct {
	Asset            oc.SymbolId `json:"asset,omitempty"`
	NeedBtcValuation *bool       `json:"needBtcValuation,omitempty"`
}

type UserAssetInfo struct {
	Asset        oc.SymbolId `json:"asset"`
	Free         oc.Amount   `json:"free"`
	Locked       oc.Amount   `json:"locked"`
	Freeze       oc.Amount   `json:"freeze"`
	Withdrawing  oc.Amount   `json:"withdrawing"`
	Ipoable      oc.Amount   `json:"ipoable"`
	BtcValuation oc.Amount   `json:"btcValuation"`
}

// https://developers.binance.com/docs/wallet/asset/user-assets
func (c *Client) GetUserAsset(args *UserAssetRequest) ([]UserAssetInfo, error) {
	var response []UserAssetInfo
	query := url.Values{}

	if args != nil {
		if args.Asset != "" {
			query.Set("asset", string(args.Asset))
		}
		if args.NeedBtcValuation != nil {
			if *args.NeedBtcValuation {
				query.Set("needBtcValuation", "true")
			} else {
				query.Set("needBtcValuation", "false")
			}
		}
	}
	query.Set("timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/1000))

	_, err := c.Request("POST", "/sapi/v3/asset/getUserAsset", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return response, nil
}
