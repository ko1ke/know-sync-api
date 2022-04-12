package pagination_utils

type Pagination struct {
	ItemsCount  int  `json:"itemsCount"`
	CountAll    int  `json:"countAll"`
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	IsFirst     bool `json:"isFirst"`
	IsLast      bool `json:"isLast"`
}

func NewPagination(page, limit, itemsCount, countAll int) *Pagination {
	totalPages := (itemsCount)/limit + 1
	var isFirst bool
	var isLast bool

	if page <= 1 {
		isFirst = true
	} else {
		isFirst = false
	}
	if totalPages <= page {
		isLast = true
	} else {
		isLast = false
	}

	return &Pagination{
		ItemsCount:  itemsCount,
		CountAll:    countAll,
		CurrentPage: page,
		TotalPages:  totalPages,
		IsFirst:     isFirst,
		IsLast:      isLast,
	}
}
