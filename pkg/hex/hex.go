package hex

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

var EncodeToString = hex.EncodeToString
var DecodeString = hex.DecodeString

type Hex []byte

func (h Hex) String() string {
	return hex.EncodeToString(h)
}

func (h Hex) Bytes() []byte {
	return []byte(h)
}

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hex) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	bz, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		return err
	}
	*h = bz
	return nil
}
