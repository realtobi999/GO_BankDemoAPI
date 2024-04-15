package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func CalculateAge(birthday time.Time) int {
	today := time.Now()
	age := today.Year() - birthday.Year()
	if today.Month() < birthday.Month() || (today.Month() == birthday.Month() && today.Day() < birthday.Day()) {
		age--
	}
	return age
}

func ParseLimitOffsetParams(r *http.Request) (int, int, error) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var limit, offset int
	var err error

	// Parse limit parameter
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, err
		}
	}

	// Parse offset parameter
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, err
		}
	}

	// Set default limit
	if limit == 0 {
		limit = 50
	}

	return limit, offset, nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}