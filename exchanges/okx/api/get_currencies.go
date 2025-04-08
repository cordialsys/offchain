package api

import (
	"fmt"

	oc "github.com/cordialsys/offchain"
	rebal "github.com/cordialsys/offchain"
)

type GetCurrenciesResponse = Response[[]CurrencyInfo]

type CurrencyInfo struct {
	BurningFeeRate       string                `json:"burningFeeRate"`
	CanDep               bool                  `json:"canDep"`
	CanInternal          bool                  `json:"canInternal"`
	CanWd                bool                  `json:"canWd"`
	Currency             oc.SymbolId           `json:"ccy"`
	Chain                SymbolAndChain        `json:"chain"`
	ContractAddress      rebal.ContractAddress `json:"ctAddr"`
	DepEstOpenTime       string                `json:"depEstOpenTime"`
	DepQuotaFixed        string                `json:"depQuotaFixed"`
	DepQuoteDailyLayer2  string                `json:"depQuoteDailyLayer2"`
	Fee                  string                `json:"fee"`
	LogoLink             string                `json:"logoLink"`
	MainNet              bool                  `json:"mainNet"`
	MaxFee               string                `json:"maxFee"`
	MaxFeeForCtAddr      string                `json:"maxFeeForCtAddr"`
	MaxWithdraw          string                `json:"maxWd"`
	MinDeposit           string                `json:"minDep"`
	MinDepArrivalConfirm string                `json:"minDepArrivalConfirm"`
	MinFee               string                `json:"minFee"`
	MinFeeForCtAddr      string                `json:"minFeeForCtAddr"`
	MinInternal          string                `json:"minInternal"`
	MinWithdraw          string                `json:"minWd"`
	MinWdUnlockConfirm   string                `json:"minWdUnlockConfirm"`
	Name                 string                `json:"name"`
	NeedTag              bool                  `json:"needTag"`
	UsedDepQuotaFixed    string                `json:"usedDepQuotaFixed"`
	UsedWdQuota          string                `json:"usedWdQuota"`
	WdEstOpenTime        string                `json:"wdEstOpenTime"`
	WithdrawQuota        string                `json:"wdQuota"`
	WithdrawTickSize     string                `json:"wdTickSz"`
}

// https://www.okx.com/docs-v5/en/#funding-account-rest-api-get-currencies
func (c *Client) GetCurrencies() (*GetCurrenciesResponse, error) {
	var response GetCurrenciesResponse
	_, err := c.Request("GET", "/api/v5/asset/currencies", nil, &response, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get assets: %w", err)
	}
	return &response, nil
}
