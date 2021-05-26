package utils

import (
	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/snippet"
)

// =====================================================================

// SplitSnippetSlice [Tasks 2.1 and 3.3.1]
// Разделяет слайс сниппетов (snippetSlice) на батчи фиксированного размера (batchSize)
// Ошибка, если размер батча нулевой.
func SplitSnippetSlice(snippetSlice *snippet.Snippets, batchSize uint) ([]snippet.Snippets, error) {
	if batchSize == 0 {
		return nil, errors.New("in func SplitSnippetSlice: invalid batchSize!")
	}

	// Обрабатывается отдельно, чтобы не создавать пустой батч на "несуществующий" слайс
	if snippetSlice == nil {
		return nil, nil
	}

	// Длина результируещего массива равна отношению длины входного слайса к размеру батча, округлённому вверх: (x - 1)/y + 1.
	res := make([]snippet.Snippets, (len(*snippetSlice)-1)/int(batchSize)+1)
	// lastBatch - последний на данный момент батч. Его размер точно не известен, т.к. зависит от размера слайса, но известен максимальный размер.
	// Объявлен тут, чтобы не делать этого на каждой итерации цикла.
	lastBatch := make(snippet.Snippets, 0, batchSize)

	for idx, snipp := range *snippetSlice {
		if idx > 0 && idx%int(batchSize) == 0 {
			res[idx/int(batchSize)-1] = lastBatch
			lastBatch = make(snippet.Snippets, 0, batchSize)
		}

		lastBatch = append(lastBatch, snipp)
	}

	if len(lastBatch) > 0 {
		res[len(res)-1] = lastBatch
	}

	return res, nil
}

// =====================================================================

// =====================================================================

// ReverseSnippetMap [Task 2.2].
// Конвертирует отображение uint64<->*Snippet в отображение *Snippet<->uin64.
// Если в исходном отображении значения дублируются вызывается паника.
func ReverseSnippetMap(snippetMap *snippet.SnippetMap) map[*snippet.Snippet]uint64 {
	res := make(map[*snippet.Snippet]uint64, len(*snippetMap))

	for key, value := range *snippetMap {
		if _, found := res[value]; found {
			panic("key duplicate in reverse map!")
		}

		res[value] = key
	}

	return res
}

// =====================================================================

// =====================================================================

// FilterSnippetSlice [Task 2.3]
// Фильтрует слайс сниппетов (snippetSlice) по фильтру (filter).
// Возвращает слайс сниппетов, в котором отсутствуют элементы из фильтра.
// filter сделан в виде переменного числа аргументов, чтобы можно было передавать array.
func FilterSnippetSlice(snippetSlice *snippet.Snippets, filter ...*snippet.Snippet) snippet.Snippets {
	if len(filter) == 0 {
		return *snippetSlice
	}

	// Заранее точную длину не знаем, т.к. в слайсе элементы не уникальны, а в фильтре могут быть элементы, отсутствующие в слайсе.
	// Но результат будет не длиннее входного слайса.
	res := make(snippet.Snippets, 0, len(*snippetSlice))

	filtered := func(snippetPtr *snippet.Snippet) bool {
		for _, filterItem := range filter {
			if snippetPtr == filterItem || *snippetPtr == *filterItem {
				return true
			}
		}
		return false
	}

	for _, snippet := range *snippetSlice {
		if !filtered(snippet) {
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
func SliceToMap(slice *snippet.Snippets) (snippet.SnippetMap, error) {
	res := make(snippet.SnippetMap, len(*slice))

	for _, value := range *slice {
		if _, found := res[value.Id]; found {
			return nil, errors.New("snippet Id duplicate!")
		}
		res[value.Id] = value
	}

	return res, nil
}

// =====================================================================
