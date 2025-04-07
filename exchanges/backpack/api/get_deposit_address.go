package api

import (
	"fmt"
	"net/url"

	oc "github.com/cordialsys/offchain"
)

type DepositAddressRequest struct {
	Blockchain oc.NetworkId `json:"blockchain"`
}

// DepositAddressResponse represents the response from the deposit address endpoint
type DepositAddressResponse struct {
	Address string `json:"address"`
}

func (c *Client) GetDepositAddress(req *DepositAddressRequest) (*DepositAddressResponse, error) {
	if req == nil || req.Blockchain == "" {
		return nil, fmt.Errorf("blockchain is required")
	}

	query := url.Values{}
	query.Set("blockchain", string(req.Blockchain))

	var response DepositAddressResponse
	_, err := c.Request("GET", "/wapi/v1/capital/deposit/address", "depositAddressQuery", nil, &response, query)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
