package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const MaxBodySize = 1 << 20 // 1 MB

func ParseJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		fmt.Println("error", err)
		return err
	}
	return nil
}
