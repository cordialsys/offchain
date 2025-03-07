package api

import "fmt"

type Response[T any] struct {
	RetCode    RetCode                `json:"retCode"`
	RetMsg     string                 `json:"retMsg"`
	Result     T                      `json:"result"`
	RetExtInfo map[string]interface{} `json:"retExtInfo"`
	Time       int64                  `json:"time"`
}

type RetCode int

func (r RetCode) String() string {
	switch r {
	case 0:
		return "success"
	case 10000:
		return "server timeout"
	case 10001:
		return "request parameter error"
	case 10002:
		return "request time exceeds the time window range"
	case 10003:
		return "API key is invalid. Check whether the key and domain are matched"
	case 33004:
		return "your API key has expired"
	case 10004:
		return "error sign, please check your signature generation algorithm"
	case 10005:
		return "permission denied, please check your API key permissions"
	case 10006:
		return "too many visits, exceeded the API Rate Limit"
	case 10007:
		return "user authentication failed"
	case 10008:
		return "common banned, please check your account mode"
	case 10009:
		return "IP has been banned"
	case 10010:
		return "unmatched IP, please check your API key's bound IP addresses"
	case 10014:
		return "invalid duplicate request"
	case 10016:
		return "server error"
	case 10017:
		return "route not found"
	case 10018:
		return "exceeded the IP Rate Limit"
	case 10024:
		return "compliance rules triggered"
	case 10027:
		return "transactions are banned"
	case 10028:
		return "the API can only be accessed by unified account users"
	case 10029:
		return "the requested symbol is invalid, please check symbol whitelist"
	default:
		return fmt.Sprintf("unknown error code (%d)", r)
	}
}
