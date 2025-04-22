package api

import (
	"net/http"
	"strconv"
	"strings"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// QueryParams holds common query parameters
type QueryParams struct {
	Page    int
	PerPage int
	Sort    string
	Order   Order
	Filter  map[string]string
}

// DefaultParams returns default query parameters
func DefaultParams() *QueryParams {
	return &QueryParams{
		Page:    1,
		PerPage: 10,
		Sort:    "id",
		Order:   OrderAsc,
		Filter:  make(map[string]string),
	}
}

// ParseQueryParams parses query parameters from a request
func ParseQueryParams(r *http.Request) *QueryParams {
	param := DefaultParams()

	getQuery := r.URL.Query().Get

	page := getQuery("page")
	// Parse pagination parameters
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			param.Page = p
		}
	}

	perPage := getQuery("perPage")
	if perPage != "" {
		if p, err := strconv.Atoi(perPage); err == nil && p > 0 {
			param.PerPage = p
		}
	}

	sort := getQuery("sort")
	if sort != "" {
		param.Sort = sort
	}

	order := getQuery("order")
	if order == "asc" || order == "desc" {
		param.Order = Order(order)
	}

	for key, values := range r.URL.Query() {
		_, ok := param.Filter[key]
		if !ok {
			length := len(values)
			if length > 1 {
				param.Filter[key] = strings.Join(values, ",")
			} else if length == 1 {
				param.Filter[key] = values[0]
			}
		}
	}

	return param
}

// GetOffset calculates the offset for pagination
func (p *QueryParams) GetOffset() int {
	return (p.Page - 1) * p.PerPage
}

// GetLimit returns the limit for pagination
func (p *QueryParams) GetLimit() int {
	return p.PerPage
}
