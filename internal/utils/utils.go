package utils

import (
	"errors"

	"github.com/ozoncp/ocp-snippet-api/internal/snippet"
)

// =====================================================================
// 2.1 (3.3.1) Split slice on fixed size batches:
//type SnippetSlice = []*snippet.Snippet
type SnippetSliceBatched = []snippet.SnippetSlice

func SplitSnippetSlice(snippetSlice *snippet.SnippetSlice, batchSize uint) (SnippetSliceBatched, error) {
	if batchSize <= 0 {
		return SnippetSliceBatched{}, errors.New("in func SplitSnippetSlice: invalid batchSize!")
	}

	if len(*snippetSlice) == 0 {
		return SnippetSliceBatched{}, nil
	}

	// Длина результируещего массива равна отношению длины входного слайса к размеру батча, округлённому вверх: (x - 1)/y + 1.
	res := make(SnippetSliceBatched, (len(*snippetSlice)-1)/int(batchSize)+1)
	// lastBatch - последний на данный момент батч. Его размер точно не известен, т.к. зависит от размера слайса, но известен максимальный размер.
	lastBatch := make(snippet.SnippetSlice, 0, batchSize)

	for idx, snipp := range *snippetSlice {
		if idx > 0 && idx%int(batchSize) == 0 {
			res[idx/int(batchSize)-1] = lastBatch
			lastBatch = make(snippet.SnippetSlice, 0, batchSize)
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
// 2.2 Reversing map:
type ReversedSnippetMap = map[*snippet.Snippet]uint64

func ReverseSnippetMap(snippetMap *snippet.SnippetMap) ReversedSnippetMap {
	res := make(ReversedSnippetMap, len(*snippetMap))

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
// 2.3 FilterSnippetSlice
//  snippetSlice - слайс сниппетов, который надо отфильтровать
//  filter - элементы, которые надо удалить из исходного слайса. Сделан в виде переменного числа аргументов, чтобы можно было передавать array.
func FilterSnippetSlice(snippetSlice *snippet.SnippetSlice, filter ...*snippet.Snippet) snippet.SnippetSlice {
	if len(filter) == 0 {
		return *snippetSlice
	}

	// Заранее точную длину не знаем, т.к. в слайсе элементы не уникальны, а в фильтре могут быть элементы, отсутствующие в слайсе.
	// Но результат будет не длиннее входного слайса.
	res := make(snippet.SnippetSlice, 0, len(*snippetSlice))

	filtered := func(snippetPtr *snippet.Snippet) bool {
		for _, filterItem := range filter {
			if snippet.CompareSnippets(snippetPtr, filterItem) { // TO BE FIXED: исправить/убрать сравнение по значению
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
// 3.3.2 SnippetSlice to SnippetMap
func SliceToMap(slice *snippet.SnippetSlice) (snippet.SnippetMap, error) {
	res := make(snippet.SnippetMap, len(*slice))

	for _, value := range *slice {
		if _, found := res[value.UserId]; found {
			return nil, errors.New("snippet UserId duplicate!")
		}
		res[value.UserId] = value
	}

	return res, nil
}

// =====================================================================
