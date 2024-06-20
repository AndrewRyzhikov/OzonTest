package inmemory

import (
	"sort"
	"time"

	"OzonTest/internal/entity"
)

type Timestamped interface {
	GetTimestamp() time.Time
}

func partition[T Timestamped](v []T, l, r int) int {
	pivot := v[(l+r)/2]
	for l <= r {
		for v[l].GetTimestamp().Before(pivot.GetTimestamp()) {
			l++
		}
		for v[r].GetTimestamp().After(pivot.GetTimestamp()) {
			r--
		}
		if l >= r {
			break
		}
		v[l], v[r] = v[r], v[l]
		l++
		r--
	}
	return r
}

func findKthSmallest[T Timestamped](v []T, l, r, k int) {
	if l >= r {
		return
	}

	m := partition[T](v, l, r)
	if k > m {
		findKthSmallest(v, m+1, r, k)
	} else {
		findKthSmallest(v, l, m, k)
	}
}

// SortWithPagination Instead of the classic sorting and taking a slice, this algorithm allows you to find the limit + offset of the earliest elements
// and sort these elements already. Thanks to this approach, the asymptotic from O(n * log(n)) is reduced to O((limit + offset) * log(limit + offset))
func SortWithPagination[T Timestamped](items []T, pagination entity.Pagination) []T {
	offset, limit := 0, len(items)

	if pagination.Offset != nil {
		offset = *pagination.Offset
	}

	if pagination.Limit != nil {
		limit = *pagination.Limit
	}

	copyItems := make([]T, len(items))
	copy(copyItems, items)

	size := min(limit+offset, len(copyItems))
	findKthSmallest(copyItems, 0, len(copyItems)-1, size-1)
	sortingItems := make([]T, size)
	copy(sortingItems, copyItems[:size])

	sort.Slice(sortingItems, func(i, j int) bool {
		return sortingItems[i].GetTimestamp().Before(sortingItems[j].GetTimestamp())
	})

	result := make([]T, min(limit, size))
	for i := 0; i < len(result) && i+offset < size; i++ {
		result[i] = sortingItems[i+offset]
	}

	return result
}
