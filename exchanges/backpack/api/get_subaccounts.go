package api

type SubaccountMapping struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type SubaccountResponse []SubaccountMapping

// No documentation for this endpoint, had to snoop it on the browser
func (c *Client) GetSubaccounts() (SubaccountResponse, error) {
	var response SubaccountResponse
	_, err := c.Request("GET", "/wapi/v1/subaccount", "subaccountQueryAll", nil, &response, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}
