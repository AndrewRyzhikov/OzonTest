package inmemory

import (
	"OzonTest/internal/entity"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSortPosts(t *testing.T) {
	posts := []entity.Post{
		{PostID: 1, Timestamp: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
		{PostID: 2, Timestamp: time.Date(2022, 12, 31, 9, 0, 0, 0, time.UTC)},
		{PostID: 3, Timestamp: time.Date(2023, 1, 2, 8, 0, 0, 0, time.UTC)},
		{PostID: 4, Timestamp: time.Date(2022, 12, 30, 7, 0, 0, 0, time.UTC)},
		{PostID: 5, Timestamp: time.Date(2023, 1, 3, 6, 0, 0, 0, time.UTC)},
	}

	offset := 0
	limit := 3
	pagination := entity.Pagination{Offset: &offset, Limit: &limit}

	expected := []entity.Post{
		{PostID: 4, Timestamp: time.Date(2022, 12, 30, 7, 0, 0, 0, time.UTC)},
		{PostID: 2, Timestamp: time.Date(2022, 12, 31, 9, 0, 0, 0, time.UTC)},
		{PostID: 3, Timestamp: time.Date(2023, 1, 2, 8, 0, 0, 0, time.UTC)},
	}

	result := sortPosts(posts, pagination)

	assert.Equal(t, len(expected), len(result))
	for i := range expected {
		assert.Equal(t, expected[i].PostID, result[i].PostID)
		assert.Equal(t, expected[i].Timestamp, result[i].Timestamp)
	}
}

func TestSortPostsWithOffset(t *testing.T) {
	// Создаем тестовые данные
	posts := []entity.Post{
		{PostID: 1, Timestamp: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)},
		{PostID: 2, Timestamp: time.Date(2022, 12, 31, 9, 0, 0, 0, time.UTC)},
		{PostID: 3, Timestamp: time.Date(2023, 1, 2, 8, 0, 0, 0, time.UTC)},
		{PostID: 4, Timestamp: time.Date(2022, 12, 30, 7, 0, 0, 0, time.UTC)},
		{PostID: 5, Timestamp: time.Date(2023, 1, 3, 6, 0, 0, 0, time.UTC)},
	}

	offset := 2
	limit := 2
	pagination := entity.Pagination{Offset: &offset, Limit: &limit}

	// Ожидаемый результат
	expected := []entity.Post{
		{PostID: 3, Timestamp: time.Date(2023, 1, 2, 8, 0, 0, 0, time.UTC)},
		{PostID: 5, Timestamp: time.Date(2023, 1, 3, 6, 0, 0, 0, time.UTC)},
	}

	// Вызов тестируемой функции
	result := sortPosts(posts, pagination)

	// Сравнение результата с ожидаемым значением
	assert.Equal(t, len(expected), len(result))
	for i := range expected {
		assert.Equal(t, expected[i].PostID, result[i].PostID)
		assert.Equal(t, expected[i].Timestamp, result[i].Timestamp)
	}
}
