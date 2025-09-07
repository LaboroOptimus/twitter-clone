package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

const MaxBodySize = 1 << 20 // 1 MB

// dublicate
func ParseJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

func ParseAndValidateJSON(
	w http.ResponseWriter,
	r *http.Request,
	validator *validator.Validate,
	dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySize)
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		fmt.Println("return 1", err.Error())
		return err
	}

	if err := validator.Struct(dst); err != nil {
		fmt.Println("return 1")
		return err
	}

	return nil
}
