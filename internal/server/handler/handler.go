package handler

type ListResponse[T any] struct {
	Items    []T  `json:"items"`
	Continue bool `json:"continue,omitempty"`
	Page     int  `json:"page,omitempty"`
	PerPage  int  `json:"per_page,omitempty"`
}
