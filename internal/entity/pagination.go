package entity

type Pagination struct {
	Limit  *int
	Offset *int
}

func NewPagination(limit *int, offset *int) *Pagination {
	return &Pagination{Limit: limit, Offset: offset}
}
