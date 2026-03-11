package server

import "encoding/json"

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
