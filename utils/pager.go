package utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	DefaultLimit = 10
	MaxLimit = 20
)

type Pager struct {
	Skip, Limit int64
	Prev, Next *string
}

func GetPager(r *http.Request) (*Pager, error) {
	page, skip, limit, parseErr := parseSkipLimit(r, DefaultLimit, MaxLimit)
	if parseErr != nil {
		return nil, parseErr
	}
	prev, next := getNextPreviousPager(r.URL.Path, page, limit)
	return &Pager{
		Skip: skip,
		Limit: limit,
		Prev: prev,
		Next: next,
	}, nil
}

func parseSkipLimit(r *http.Request, def, max int) (int64, int64, int64, error) {
	q := r.URL.Query()
	var page, skip, limit int64 = 1, 0, int64(def)
	pageQ := q.Get("page")
	if pageQ != "" {
		pageInt, err := strconv.Atoi(pageQ)
		if err != nil {
			return page, skip, limit, errors.New("failed to parse page")
		}
		page = int64(pageInt)
	}
	limitQ := q.Get("limit")

	if limitQ != "" {
		limitInt, err := strconv.Atoi(limitQ)
		if err != nil {
			return page, skip, limit, errors.New("failed to parse limit")
		}
		limit = int64(limitInt)
	}

	if limit > int64(max) {
		limit = int64(max)
	}
	skip = limit * (page - 1)
	if skip < 0 {
		skip = 0
	}
	return page, skip, limit, nil
}

func getNextPreviousPager(path string, page, limit int64) (*string, *string) {
	var previous, next string
	if page-1 > 0 {
		previous = fmt.Sprintf("%s?limit=%d&page=%d", path, limit, page-1)
	}
	next = fmt.Sprintf("%s?limit=%d&page=%d", path, limit, page+1)
	return &previous, &next
}

