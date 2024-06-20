package inmemory

import (
	"OzonTest/internal/entity"
	"fmt"
	"sort"
)

func partition(v []entity.Post, l, r int) int {
	pivot := v[(l+r)/2]
	for l <= r {
		for v[l].Timestamp.Before(pivot.Timestamp) {
			l++
		}
		for v[r].Timestamp.After(pivot.Timestamp) {
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

func findKthSmallest(posts []entity.Post, l, r, k int) {
	if l == r {
		return
	}

	m := partition(posts, l, r)
	if k > m {
		findKthSmallest(posts, m+1, r, k)
	} else {
		findKthSmallest(posts, l, m, k)
	}
}

func sortPosts(posts []entity.Post, pagination entity.Pagination) []entity.Post {
	p := make([]entity.Post, 0)
	for i := *pagination.Offset; i < len(posts); i++ {
		p = append(p, posts[i])
	}

	findKthSmallest(p, 0, len(p)-1, *pagination.Limit-1)

	result := make([]entity.Post, *pagination.Limit)
	for i := 0; i < *pagination.Limit; i++ {
		result[i] = posts[i]
	}

	fmt.Println(result)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})
	fmt.Println(result)

	return result
}
