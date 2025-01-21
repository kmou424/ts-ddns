package util

import (
	"encoding/json"
	"net/http"
)

func MarshalResp[T any](resp *http.Response) (*T, error) {
	val := new(T)
	err := json.NewDecoder(resp.Body).Decode(val)
	if err != nil {
		panic(err)
	}
	return val, nil
}
