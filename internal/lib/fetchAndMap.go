package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

/* FetchAndMap sends and recieves http request and maps it to T via JSON */
func FetchAndMap[T any](fetch func() (*http.Response, error)) (T, error) {
	var ret T

	res, err := fetch()
	if err != nil {
		return ret, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
