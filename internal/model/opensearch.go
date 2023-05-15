package model

import (
	"fmt"
	"regexp"
	"strings"
)

type OSPaginationRequest struct {
	From           int64             `json:"from"`
	Size           int64             `json:"size"`
	TrackTotalHits bool              `json:"track_total_hits"`
	Query          Query             `json:"query"`
	Sort           []map[string]Sort `json:"sort"`
}

type Query struct {
	Bool Bool `json:"bool"`
}

type Bool struct {
	MustNot []*MustNot `json:"must_not,omitempty"`
	Must    *Must      `json:"must,omitempty"`
	Filter  []Filter   `json:"filter,omitempty"`
}

type Filter struct {
	Term map[string]string `json:"term"`
}

type Must struct {
	MultiMatch MultiMatch `json:"multi_match"`
}

type MultiMatch struct {
	Query              string   `json:"query"`
	Analyzer           string   `json:"analyzer"`
	Fields             []string `json:"fields"`
	MinimumShouldMatch string   `json:"minimum_should_match"`
}

type MustNot struct {
	Exists Exists `json:"exists"`
}

type Exists struct {
	Field string `json:"field"`
}

// -created_at = created_at desc.
// +created_at = created_at asc.
type Sort struct {
	Order string `json:"order"`
}

func (m *OSPaginationRequest) ParseSort(req *PaginationPayload) {
	sortReq := make([]map[string]Sort, 0)
	for _, sort := range req.Sort {
		condition := "desc"
		if strings.Contains(sort, "+") {
			condition = "asc"
		}
		re := regexp.MustCompile(`[+-]`)
		sort = re.ReplaceAllString(sort, "")
		sort := fmt.Sprintf("%s.keyword", sort)
		sortReq = append(sortReq, map[string]Sort{
			sort: {
				Order: condition,
			},
		})
	}
	m.Sort = sortReq
}

type OSPaginationResponse[T any] struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source *T `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (m *OSPaginationResponse[t]) GetCount() int64 {
	return m.Hits.Total.Value
}

func (m *OSPaginationResponse[T]) GetItems() []*T {
	results := make([]*T, 0)
	for _, item := range m.Hits.Hits {
		results = append(results, item.Source)
	}
	return results
}

type IndexModel interface {
	GetID() string
	ToDoc() any
}
