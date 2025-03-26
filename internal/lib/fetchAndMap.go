package lib

import (
	"encoding/json"
	"io"
	"net/http"
)

/* Mapping module */

/* FetchAndMap sends and receives http request and maps it to T via JSON */
func FetchAndMap[T any](fetch func() (*http.Response, error)) (T, *http.Response, error) {
	var ret T

	res, err := fetch()
	if err != nil {
		return ret, res, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ret, res, err
	}

	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, res, err
	}

	return ret, res, nil
}

func Map[T any](d []byte) (T, error) {
	var t T
	if err := json.Unmarshal(d, &t); err != nil {
		return t, err
	}
	return t, nil
}
