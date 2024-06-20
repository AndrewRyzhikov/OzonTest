package inmemory

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"OzonTest/internal/entity"
)

func generatePosts(length int) []*entity.Post {
	slice := make([]*entity.Post, length)

	for i := 0; i < length; i++ {
		slice[i] = &entity.Post{
			ID:        i + 1,
			Timestamp: time.Now().Add(time.Duration(i) * time.Hour * -1),
		}
	}

	return slice
}

type testCase struct {
	name                  string
	offset, limit, length int
}

func generateTestCases(num int) []testCase {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	testCases := make([]testCase, num)
	for i := 0; i < num; i++ {
		testCases[i] = testCase{
			name:   fmt.Sprintf("Test case %d", i+1),
			offset: r.Intn(1000),
			limit:  r.Intn(1000),
			length: r.Intn(1000),
		}
	}
	return testCases
}

func ClassicSort(generation []*entity.Post, offset, limit, length int) []*entity.Post {
	sort.Slice(generation, func(i, j int) bool {
		return generation[i].Timestamp.Before(generation[j].Timestamp)
	})

	expected := make([]*entity.Post, min(limit, length))
	for i := 0; i < len(expected) && i+offset < length; i++ {
		expected[i] = generation[i+offset]
	}

	return expected
}

func ClassicSortAndMySortPosts(offset, limit, length int) ([]*entity.Post, []*entity.Post) {
	generation := generatePosts(length)
	pagination := entity.Pagination{Offset: &offset, Limit: &limit}

	posts1 := make([]*entity.Post, length)
	copy(posts1, generation)
	result := SortWithPagination[*entity.Post](posts1, pagination)

	posts2 := make([]*entity.Post, length)
	copy(posts2, generation)
	expected := ClassicSort(posts2, offset, limit, length)

	return expected, result
}

func TestSortPosts(t *testing.T) {
	tests := generateTestCases(1000)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected, result := ClassicSortAndMySortPosts(tt.offset, tt.limit, tt.length)

			assert.Equal(t, len(expected), len(result))
			for i := range expected {
				if expected[i] == nil && result[i] == nil {
					continue
				}
				assert.Equal(t, expected[i].ID, result[i].ID)
				assert.Equal(t, expected[i].Timestamp, result[i].Timestamp)
			}
		})
	}
}
