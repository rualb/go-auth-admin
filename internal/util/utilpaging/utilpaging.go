package utilpaging

type PagingInputDTO struct {
	Page   int    `query:"page"`
	Sort   string `query:"sort"`
	Search string `query:"search"`
	Cursor string `query:"cursor"`
	// Filter   map[string][]string //  c.QueryParams() ref
	Limit int `query:"limit"`
}

// func (x *PagingInputDTO) GetFilter(code string) string {
// 	r := x.Filter[code]
// 	if len(r) > 0 {
// 		return r[0]
// 	}
// 	return ""
// }

func (x *PagingInputDTO) Info(totalCount int) *PagingInfo {

	res := &PagingInfo{}
	res.Fill(totalCount, x.Limit, x.Page)
	return res
}

// PagingInfo holds pagination information
type PagingInfo struct {
	Offset          int  // Offset
	Limit           int  // Number of items per page
	Page            int  // Current page number
	PageCount       int  // Total number of pages
	TotalCount      int  // Total number of items
	HasNextPage     bool // Whether there is a next page
	HasPreviousPage bool // Whether there is a previous page
	IsFirstPage     bool // Whether the current page is the first page
	IsLastPage      bool // Whether the current page is the last page
}

func (x *PagingInfo) Fill(totalCount int, limit int, page int) {
	// Ensure Page and Limit have sensible defaults
	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10 // Default page size
	}

	if limit > 1000 {
		limit = 1000 // Max page size
	}

	// Calculate total page count
	pageCount := (totalCount + limit - 1) / limit

	// if pageCount*limit < totalCount {
	// 	pageCount++
	// }

	// Adjust page number if it exceeds the total count of pages
	if page > pageCount {
		page = pageCount
	}

	// Calculate offset
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	// Create and return PaginationInfo

	x.Offset = offset
	x.Limit = limit
	x.Page = page
	x.PageCount = pageCount
	x.TotalCount = totalCount
	x.HasNextPage = x.Page < pageCount
	x.HasPreviousPage = x.Page > 1
	x.IsFirstPage = x.Page == 1
	x.IsLastPage = x.Page == pageCount

}

// generateNavPages creates an array of page numbers for navigation
func NavPages(currentPage, totalPages int) []int {
	var navPages []int
	for i := -2; i <= 2; i++ {
		page := currentPage + i
		if page >= 1 && page <= totalPages {
			navPages = append(navPages, page)
		}
	}
	// Include the first and last page if not already included
	if currentPage > 3 {
		navPages = append(navPages, 1)
	}
	if currentPage < totalPages-2 {
		navPages = append(navPages, totalPages)
	}
	return unique(navPages)
}

// unique removes duplicate integers from a slice
func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	var list []int
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type PagingOutputDTO[T any] struct {
	Filter struct {
		Page   int    `json:"page,omitempty"`
		Sort   string `json:"sort,omitempty"`
		Search string `json:"search,omitempty"`
		Cursor string `json:"cursor,omitempty"`
		// Filter     map[string][]string `json:"filter"`
		Limit int `json:"limit,omitempty"`
	} `json:"filter,omitempty"`

	Info struct {
		PageCount  int `json:"page_count,omitempty"`
		TotalCount int `json:"total_count,omitempty"`
	} `json:"info,omitempty"`

	Data []*T `json:"data,omitempty"`
}

func (x *PagingOutputDTO[T]) Fill(filter *PagingInputDTO, info *PagingInfo) {
	x.Filter.Page = info.Page
	x.Filter.Limit = info.Limit
	x.Filter.Sort = filter.Sort
	x.Filter.Search = filter.Search

	//
	// x.Cursor = info.Cursor
	// x.Filter = filter.Filter
	//
	x.Info.PageCount = info.PageCount
	x.Info.TotalCount = info.TotalCount
}
