package dto

import "encoding/json"

func unmarshalJSONStrings(data []byte) []string {
	if data == nil {
		return []string{}
	}
	var arr []string
	_ = json.Unmarshal(data, &arr) // ignore error; if invalid, return empty slice
	return arr
}
