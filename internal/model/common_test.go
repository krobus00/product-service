package model

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name string
		want *Response
	}{
		{
			name: "success",
			want: &Response{
				Message: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultResponse(t *testing.T) {
	tests := []struct {
		name string
		want *Response
	}{
		{
			name: "success",
			want: &Response{
				Message: "OK",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_WithMessage(t *testing.T) {
	type fields struct {
		Message string
		Data    any
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "success",
			fields: fields{
				Message: "message",
			},
			args: args{
				message: "message",
			},
			want: &Response{
				Message: "message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Response{
				Message: tt.fields.Message,
				Data:    tt.fields.Data,
			}
			if got := m.WithMessage(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.WithMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_WithErrorMessage(t *testing.T) {
	type fields struct {
		Message string
		Data    any
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "success",
			fields: fields{
				Message: "sample error message",
			},
			args: args{
				err: errors.New("sample error message"),
			},
			want: &Response{
				Message: "sample error message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Response{
				Message: tt.fields.Message,
				Data:    tt.fields.Data,
			}
			if got := m.WithErrorMessage(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.WithErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResponse_WithData(t *testing.T) {
	type fields struct {
		Message string
		Data    any
	}
	type args struct {
		data any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Response
	}{
		{
			name: "success",
			fields: fields{
				Data: "sample string",
			},
			args: args{
				data: "sample string",
			},
			want: &Response{
				Data: "sample string",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Response{
				Message: tt.fields.Message,
				Data:    tt.fields.Data,
			}
			if got := m.WithData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Response.WithData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPaginationPayload(t *testing.T) {
	tests := []struct {
		name string
		want *PaginationPayload
	}{
		{
			name: "success",
			want: &PaginationPayload{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaginationPayload(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaginationPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationPayload_Sanitize(t *testing.T) {
	type fields struct {
		Search string
		Sort   []string
		Limit  int
		Page   int
	}
	tests := []struct {
		name   string
		fields fields
		want   *PaginationPayload
	}{
		{
			name:   "success with default value",
			fields: fields{},
			want: &PaginationPayload{
				Search: "",
				Limit:  defaultPaginationLimit,
				Page:   1,
			},
		},
		{
			name: "success with set limit to max limit",
			fields: fields{
				Limit: 9999999,
			},
			want: &PaginationPayload{
				Search: "",
				Limit:  maxPaginationLimit,
				Page:   1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationPayload{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				Limit:  tt.fields.Limit,
				Page:   tt.fields.Page,
			}
			if got := m.Sanitize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationPayload.Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationPayload_WithSearch(t *testing.T) {
	type fields struct {
		Search string
		Sort   []string
		Limit  int
		Page   int
	}
	type args struct {
		search string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationPayload
	}{
		{
			name: "success",
			fields: fields{
				Search: "query",
			},
			args: args{
				search: "query",
			},
			want: &PaginationPayload{
				Search: "query",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationPayload{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				Limit:  tt.fields.Limit,
				Page:   tt.fields.Page,
			}
			if got := m.WithSearch(tt.args.search); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationPayload.WithSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationPayload_WithSort(t *testing.T) {
	type fields struct {
		Search string
		Sort   []string
		Limit  int
		Page   int
	}
	type args struct {
		sort []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationPayload
	}{
		{
			name: "success",
			fields: fields{
				Sort: []string{"+created_at"},
			},
			args: args{
				sort: []string{"+created_at"},
			},
			want: &PaginationPayload{
				Sort: []string{"+created_at"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationPayload{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				Limit:  tt.fields.Limit,
				Page:   tt.fields.Page,
			}
			if got := m.WithSort(tt.args.sort); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationPayload.WithSort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationPayload_WithPage(t *testing.T) {
	type fields struct {
		Search string
		Sort   []string
		Limit  int
		Page   int
	}
	type args struct {
		page int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationPayload
	}{
		{
			name: "success",
			fields: fields{
				Page: 17,
			},
			args: args{
				page: 17,
			},
			want: &PaginationPayload{
				Page: 17,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationPayload{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				Limit:  tt.fields.Limit,
				Page:   tt.fields.Page,
			}
			if got := m.WithPage(tt.args.page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationPayload.WithPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationPayload_WithLimit(t *testing.T) {
	type fields struct {
		Search string
		Sort   []string
		Limit  int
		Page   int
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationPayload
	}{
		{
			name: "success",
			fields: fields{
				Limit: 10,
			},
			args: args{
				limit: 10,
			},
			want: &PaginationPayload{
				Limit: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationPayload{
				Search: tt.fields.Search,
				Sort:   tt.fields.Sort,
				Limit:  tt.fields.Limit,
				Page:   tt.fields.Page,
			}
			if got := m.WithLimit(tt.args.limit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationPayload.WithLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPaginationResponse(t *testing.T) {
	type args struct {
		request *PaginationPayload
	}
	tests := []struct {
		name string
		args args
		want *PaginationResponse
	}{
		{
			name: "success",
			args: args{
				request: NewPaginationPayload().WithLimit(10).WithPage(1),
			},
			want: NewPaginationResponse(NewPaginationPayload().WithLimit(10).WithPage(1)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaginationResponse(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaginationResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationResponse_WithCount(t *testing.T) {
	type fields struct {
		Meta    *PaginationPayload
		Count   int64
		MaxPage int64
		Items   []string
	}
	type args struct {
		count int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationResponse
	}{
		{
			name: "success",
			fields: fields{
				Count: 10,
				Items: []string{},
			},
			args: args{
				count: 10,
			},
			want: &PaginationResponse{
				Items: []string{},
				Count: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationResponse{
				Meta:    tt.fields.Meta,
				Count:   tt.fields.Count,
				MaxPage: tt.fields.MaxPage,
				Items:   tt.fields.Items,
			}
			if got := m.WithCount(tt.args.count); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationResponse.WithCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationResponse_WithItems(t *testing.T) {
	sampleItems := make([]string, 0)
	type fields struct {
		Meta    *PaginationPayload
		Count   int64
		MaxPage int64
		Items   []string
	}
	type args struct {
		items []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *PaginationResponse
	}{
		{
			name: "success",
			fields: fields{
				Items: sampleItems,
			},
			args: args{
				items: sampleItems,
			},
			want: &PaginationResponse{
				Items: sampleItems,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationResponse{
				Meta:    tt.fields.Meta,
				Count:   tt.fields.Count,
				MaxPage: tt.fields.MaxPage,
				Items:   tt.fields.Items,
			}
			if got := m.WithItems(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationResponse.WithItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationResponse_BuildResponse(t *testing.T) {
	type fields struct {
		Meta  *PaginationPayload
		Count int64
		Items []string
	}
	tests := []struct {
		name   string
		fields fields
		want   *PaginationResponse
	}{
		{
			name: "success",
			fields: fields{
				Meta:  NewPaginationPayload().WithLimit(10),
				Count: 100,
				Items: []string{},
			},
			want: &PaginationResponse{
				Meta:    NewPaginationPayload().WithLimit(10),
				Count:   100,
				MaxPage: 10,
				Items:   []string{},
			},
		},
		{
			name: "success ceil up",
			fields: fields{
				Meta:  NewPaginationPayload().WithLimit(9),
				Count: 9,
				Items: []string{},
			},
			want: &PaginationResponse{
				Meta:    NewPaginationPayload().WithLimit(9),
				Count:   9,
				MaxPage: 1,
				Items:   []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &PaginationResponse{
				Meta:  tt.fields.Meta,
				Count: tt.fields.Count,
				Items: tt.fields.Items,
			}
			if got := m.BuildResponse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PaginationResponse.BuildResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
