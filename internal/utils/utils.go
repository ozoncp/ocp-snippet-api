package utils

import (
	"errors"
)

// TO BE FIXED:
// - реализовать Snippet согласно заданию (когда оно станет понятно...)
// - вынести объявление в отдельный пакет
type Snippet struct {
	Field int
}

// =====================================================================
// 2.1 Split slice on fixed size batches:
type SnippetSlice = []*Snippet
type SnippetSliceBatched = []SnippetSlice

func SplitSnippetSlice(snippetSlice *SnippetSlice, batchSize int) (SnippetSliceBatched, error) {
	if batchSize <= 0 {
		return SnippetSliceBatched{}, errors.New("in func SplitSnippetSlice: invalid batchSize!")
	}

	if len(*snippetSlice) == 0 {
		return SnippetSliceBatched{}, nil
	}

	// Длина результируещего массива равна отношению длины входного слайса к размеру батча, округлённому вверх: (x - 1)/y + 1.
	res := make(SnippetSliceBatched, (len(*snippetSlice)-1)/batchSize+1)
	// lastBatch - последний на данный момент батч. Его размер точно не известен, т.к. зависит от размера слайса, но известен максимальный размер.
	lastBatch := make(SnippetSlice, 0, batchSize)

	for idx, snippet := range *snippetSlice {
		if idx > 0 && idx%batchSize == 0 {
			res[idx/batchSize-1] = lastBatch
			lastBatch = make(SnippetSlice, 0, batchSize)
		}

		lastBatch = append(lastBatch, snippet)
	}

	if len(lastBatch) > 0 {
		res[len(res)-1] = lastBatch
	}

	return res, nil
}

// =====================================================================

// =====================================================================
// 2.2 Reversing map:
type SnippetMap = map[int]*Snippet
type ReversedSnippetMap = map[*Snippet]int

func ReverseSnippetMap(snippetMap *SnippetMap) ReversedSnippetMap {
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
func FilterSnippetSlice(snippetSlice *SnippetSlice, filter ...*Snippet) SnippetSlice {
	if len(filter) == 0 {
		return *snippetSlice
	}

	// Заранее точную длину не знаем, т.к. в слайсе элементы не уникальны, а в фильтре могут быть элементы, отсутствующие в слайсе.
	// Но результат будет не длиннее входного слайса.
	res := make(SnippetSlice, 0, len(*snippetSlice))

	filtered := func(snippet *Snippet) bool {
		for _, filterItem := range filter {
			if snippet == filterItem || snippet.Field == filterItem.Field { // TO BE FIXED: исправить/убрать сравнение по значению
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
