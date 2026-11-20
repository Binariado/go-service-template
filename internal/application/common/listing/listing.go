package listing

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

var (
	ErrInvalidPage  = errors.New("page must be greater than 0")
	ErrInvalidLimit = errors.New("limit must be greater than 0")
	ErrInvalidSort  = errors.New("invalid sort")
	ErrInvalidOrder = errors.New("invalid order")
)

type Options struct {
	Page      int
	Limit     int
	Offset    int
	SortField string
	SortOrder string
}

func Normalize(page, limit int, sort, order string, allowedSort []string, defaultSort, defaultOrder string) (Options, error) {
	if page == 0 {
		page = DefaultPage
	}
	if page < 0 {
		return Options{}, ErrInvalidPage
	}

	if limit == 0 {
		limit = DefaultLimit
	}
	if limit < 0 {
		return Options{}, ErrInvalidLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}

	sortField := defaultSort
	if sort != "" {
		if !slices.Contains(allowedSort, sort) {
			return Options{}, fmt.Errorf("%w: %s", ErrInvalidSort, sort)
		}
		sortField = sort
	}

	sortOrder := normalizeDefaultOrder(defaultOrder)
	if order != "" {
		switch strings.ToLower(order) {
		case "asc":
			sortOrder = "ASC"
		case "desc":
			sortOrder = "DESC"
		default:
			return Options{}, fmt.Errorf("%w: %s", ErrInvalidOrder, order)
		}
	}

	return Options{
		Page:      page,
		Limit:     limit,
		Offset:    (page - 1) * limit,
		SortField: sortField,
		SortOrder: sortOrder,
	}, nil
}

func normalizeDefaultOrder(order string) string {
	if strings.EqualFold(order, "desc") {
		return "DESC"
	}
	return "ASC"
}
