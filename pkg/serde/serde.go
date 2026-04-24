package serde

import (
	"encoding/json"
	"io"
	"net/http"
)

func DecodeJSON[T any](body io.ReadCloser) (T, error) {
	var target T
	err := json.NewDecoder(body).Decode(&target)
	if err != nil {
		return target, err
	}

	return target, nil
}

func EncodeJSONHTTP[T any](data any, w http.ResponseWriter) error {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func ReadJSONHttp[T any](body io.ReadCloser, w http.ResponseWriter) (T, error) {
	target, err := DecodeJSON[T](body)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return target, err
	}

	return target, nil
}
