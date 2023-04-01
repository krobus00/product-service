package model

import (
	"errors"
	"math"

	pb "github.com/krobus00/product-service/pb/product"
)

const (
	defaultPaginationLimit = int(10)
	maxPaginationLimit     = int(20)
)

var (
	ErrUnauthorizedAccess = errors.New("unauthorized access")
)

type Response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewResponse() *Response {
	return new(Response)
}

func NewDefaultResponse() *Response {
	return &Response{
		Message: "OK",
	}
}

func (m *Response) WithMessage(message string) *Response {
	m.Message = message
	return m
}

func (m *Response) WithErrorMessage(err error) *Response {
	m.Message = err.Error()
	return m
}

func (m *Response) WithData(data any) *Response {
	m.Data = data
	return m
}

type PaginationPayload struct {
	Search         string   `json:"search" query:"search"`
	Sort           []string `json:"sort" query:"sort"`
	Limit          int      `json:"limit" query:"limit"`
	Page           int      `json:"page" query:"page"`
	IncludeDeleted bool     `json:"include_deleted" query:"includeDeleted"`
}

func NewPaginationPayloadFromProto(message *pb.PaginationRequest) *PaginationPayload {
	return &PaginationPayload{
		Search:         message.GetSearch(),
		Sort:           message.GetSort(),
		Limit:          int(message.GetLimit()),
		Page:           int(message.GetPage()),
		IncludeDeleted: message.GetIncludeDeleted(),
	}
}

func NewPaginationPayload() *PaginationPayload {
	return new(PaginationPayload)
}

func (m *PaginationPayload) ToProto() *pb.PaginationRequest {
	return &pb.PaginationRequest{
		Search:         m.Search,
		Sort:           m.Sort,
		Limit:          int64(m.Limit),
		Page:           int64(m.Page),
		IncludeDeleted: m.IncludeDeleted,
	}
}

func (m *PaginationPayload) Sanitize() *PaginationPayload {
	if m.Limit <= 0 {
		m.Limit = defaultPaginationLimit
	}
	if m.Limit > maxPaginationLimit {
		m.Limit = maxPaginationLimit
	}
	if m.Page <= 0 {
		m.Page = 1
	}
	return m
}

func (m *PaginationPayload) WithSearch(search string) *PaginationPayload {
	m.Search = search
	return m
}

func (m *PaginationPayload) WithSort(sort []string) *PaginationPayload {
	m.Sort = sort
	return m
}

func (m *PaginationPayload) WithPage(page int) *PaginationPayload {
	m.Page = page
	return m
}

func (m *PaginationPayload) WithLimit(limit int) *PaginationPayload {
	m.Limit = limit
	return m
}

type PaginationResponse struct {
	Meta    *PaginationPayload `json:"meta"`
	Count   int64              `json:"count"`
	MaxPage int64              `json:"maxPage"`
	Items   []string           `json:"items"`
}

func NewPaginationResponse(request *PaginationPayload) *PaginationResponse {
	return &PaginationResponse{
		Meta:  request,
		Items: make([]string, 0),
	}
}

func (m *PaginationResponse) ToProto() *pb.PaginationResponse {
	return &pb.PaginationResponse{
		Meta:    m.Meta.ToProto(),
		Count:   m.Count,
		MaxPage: m.MaxPage,
		Items:   m.Items,
	}
}

func (m *PaginationResponse) WithCount(count int64) *PaginationResponse {
	m.Count = count
	return m
}

func (m *PaginationResponse) WithItems(items []string) *PaginationResponse {
	if items != nil {
		m.Items = items
	}
	return m
}

func (m *PaginationResponse) BuildResponse() *PaginationResponse {
	m.MaxPage = int64(math.Ceil(float64(m.Count) / float64(m.Meta.Limit)))
	return m
}
