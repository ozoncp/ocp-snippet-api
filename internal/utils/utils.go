package utils

import (
	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/models"
)

// =====================================================================

// SplitSnippetSlice [Tasks 2.1 and 3.3.1]
// Разделяет слайс сниппетов (snippetSlice) на батчи фиксированного размера (batchSize)
// Ошибка, если размер батча нулевой.
func SplitSnippetSlice(snippets []models.Snippet, batchSize uint) ([][]models.Snippet, error) {
	if batchSize == 0 {
		return nil, errors.New("in func SplitSnippetSlice: invalid batchSize!")
	}

	// Длина результируещего массива равна отношению длины входного слайса (x) к размеру батча (y), округлённому вверх: (x - 1)/y + 1.
	res := make([][]models.Snippet, (len(snippets)-1)/int(batchSize)+1)

	for idx := range res {
		from, to := idx*int(batchSize), (idx+1)*int(batchSize)
		if to < len(snippets) {
			res[idx] = snippets[from:to]
		} else {
			res[idx] = snippets[from:]
		}
	}

	return res, nil
}

// =====================================================================

// =====================================================================

// FilterSnippetSlice [Task 2.3]
// Фильтрует слайс сниппетов (snippetSlice) по фильтру (filter).
// Возвращает слайс сниппетов, в котором отсутствуют элементы из фильтра.
// filter сделан в виде переменного числа аргументов, чтобы можно было передавать array.
func FilterSnippetSlice(snippetSlice []models.Snippet, filter ...models.Snippet) []models.Snippet {
	if len(filter) == 0 {
		return snippetSlice
	}

	// Заранее точную длину не знаем, т.к. в слайсе элементы не уникальны, а в фильтре могут быть элементы, отсутствующие в слайсе.
	// Но результат будет не длиннее входного слайса.
	res := make([]models.Snippet, 0, len(snippetSlice))

	filtered := func(snippet *models.Snippet) bool {
		for _, filterItem := range filter {
			if snippet != nil && *snippet == filterItem {
				return true
			}
		}
		return false
	}

	for _, snippet := range snippetSlice {
		if !filtered(&snippet) {
			res = append(res, snippet)
		}
	}

	return res
}

// =====================================================================

// =====================================================================
// 3.3.2 Snippets to SnippetMap

// SliceToMap [Task 3.3.2]
// Конвертация слайса сниппетов (Snippet) в отображение, ключом которого является Id, а значением - сниппет
// Если Id дублируется, возвращает ошибку.
func SliceToMap(slice []models.Snippet) (map[uint64]models.Snippet, error) {
	res := make(map[uint64]models.Snippet, len(slice))

	for _, value := range slice {
		if _, found := res[value.Id]; found {
			return nil, errors.New("snippet Id duplicate!")
		}
		res[value.Id] = value
	}

	return res, nil
}

// =====================================================================
