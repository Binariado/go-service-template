package listing

import (
	"errors"
	"testing"
)

func TestNormalizeDefaults(t *testing.T) {
	got, err := Normalize(0, 0, "", "", []string{"name"}, "name", "asc")
	if err != nil {
		t.Fatal(err)
	}
	want := Options{Page: 1, Limit: 20, Offset: 0, SortField: "name", SortOrder: "ASC"}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestNormalizeClampsLimitAndComputesOffset(t *testing.T) {
	got, err := Normalize(3, 500, "name", "desc", []string{"name"}, "name", "asc")
	if err != nil {
		t.Fatal(err)
	}
	if got.Limit != MaxLimit || got.Offset != 2*MaxLimit || got.SortOrder != "DESC" {
		t.Errorf("unexpected options %+v", got)
	}
}

func TestNormalizeErrors(t *testing.T) {
	allowed := []string{"name"}
	tests := []struct {
		name  string
		page  int
		limit int
		sort  string
		order string
		want  error
	}{
		{"page", -1, 0, "", "", ErrInvalidPage},
		{"limit", 0, -1, "", "", ErrInvalidLimit},
		{"sort", 0, 0, "id; DROP", "", ErrInvalidSort},
		{"order", 0, 0, "", "sideways", ErrInvalidOrder},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Normalize(tt.page, tt.limit, tt.sort, tt.order, allowed, "name", "asc")
			if !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}
