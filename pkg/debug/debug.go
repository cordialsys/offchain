package debug

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
