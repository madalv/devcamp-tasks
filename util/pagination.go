package util

type Pagination struct {
	page  int
	limit int
}

func NewPagination(page, limit int) *Pagination {
	return &Pagination{
		page:  page,
		limit: limit,
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.limit <= 0 {
		p.limit = 20
	}
	return p.limit
}

func (p *Pagination) GetPage() int {
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}
