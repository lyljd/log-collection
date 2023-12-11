package tool

import "encoding/json"

func IsJSON(s string) bool {
	var m map[string]any
	return json.Unmarshal([]byte(s), &m) == nil
}
