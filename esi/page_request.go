package esi

// PageRequest represents a paged request for an ESI resource.
type PageRequest struct {
	page int
	size int
}

// NewPageRequest builds a new PageRequest.
func NewPageRequest(page int, size int) PageRequest {
	if page < 0 {
		page = 0
	}

	if page == 0 {
		size = 1
	}

	if size < 1 {
		size = 1
	}

	return PageRequest{page, size}
}

// pages returns a slice containing all the page numbers based upon the starting page and size.
func (p PageRequest) pages() []int {
	pages := make([]int, p.size)
	for i := 0; i < p.size; i++ {
		pages[i] = p.page + i
	}
	return pages
}
