package utils

import (
	"testing"

	"github.com/ozoncp/ocp-snippet-api/internal/models"
)

// =====================================================================
// Testing snippets:
var (
	snippet_1     = models.Snippet{Id: 1}
	snippet_2     = models.Snippet{Id: 2}
	snippet_3     = models.Snippet{Id: 3}
	snippet_4     = models.Snippet{Id: 4}
	snippet_5     = models.Snippet{Id: 5}
	snippet_6     = models.Snippet{Id: 6}
	snippet_7     = models.Snippet{Id: 7}
	snippet_8     = models.Snippet{Id: 8}
	snippet_9     = models.Snippet{Id: 9}
	snippet_n5    = models.Snippet{Id: 15}
	snippet_19    = models.Snippet{Id: 19}
	snippet_0     = models.Snippet{Id: 0}
	snippet_empty = models.Snippet{}
)

// =====================================================================

// =====================================================================
// Compare functions:

func compareSnippetSlices(l *[]models.Snippet, r *[]models.Snippet) bool {
	if len(*l) != len(*r) {
		return false
	}

	for idx, lSnippet := range *l {
		if !(lSnippet == (*r)[idx]) {
			return false
		}
	}

	return true
}

// compareSnippetSliceBatched - функция сравнения слайсов батчей ([][]models.Snippet).
// Считает слайсы батчей равными, если их длины равны, длины батчей равны и значения снипеттов равны
func compareSnippetSliceBatched(l *[][]models.Snippet, r *[][]models.Snippet) bool {
	if len(*l) != len(*r) {
		return false
	}

	for idx, lSlice := range *l {
		if !compareSnippetSlices(&lSlice, &(*r)[idx]) {
			return false
		}
	}

	return true
}
func compareSnippetMaps(l *map[uint64]models.Snippet, r *map[uint64]models.Snippet) bool {
	if len(*l) != len(*r) {
		return false
	}

	for lKey, lValue := range *l {
		if rValue, found := (*r)[lKey]; !found || lValue != rValue {
			return false
		}
	}

	return true
}

// =====================================================================

// =====================================================================

// TestSplitSnippetSlice
func TestSplitSnippetSlice(t *testing.T) {
	type result struct {
		batchedSlice [][]models.Snippet
		errExpected  bool
	}

	testSet := []struct {
		slice     []models.Snippet
		batchSize uint
		res       result
	}{
		{
			slice:     []models.Snippet{snippet_1, snippet_2, snippet_3, snippet_4, snippet_5, snippet_6, snippet_7, snippet_8},
			batchSize: 3,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1, snippet_2, snippet_3}, {snippet_4, snippet_5, snippet_6}, {snippet_7, snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1},
			batchSize: 3,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1, snippet_2, snippet_3, snippet_4, snippet_5, snippet_6, snippet_7, snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1, snippet_2}, {snippet_3, snippet_4}, {snippet_5, snippet_6}, {snippet_7, snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1, snippet_2, snippet_3, snippet_4, snippet_5, snippet_6, snippet_7, snippet_8, snippet_9},
			batchSize: 5,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1, snippet_2, snippet_3, snippet_4, snippet_5}, {snippet_6, snippet_7, snippet_8, snippet_9}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{},
			batchSize: 5,
			res: result{
				batchedSlice: [][]models.Snippet{{}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1, snippet_3, snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1, snippet_3}, {snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1, snippet_3, snippet_8, snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_1, snippet_3, snippet_8, snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: [][]models.Snippet{{snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     []models.Snippet{snippet_1, snippet_3, snippet_8, snippet_empty},
			batchSize: 0,
			res: result{
				batchedSlice: nil,
				errExpected:  true,
			},
		},
		{
			res: result{
				batchedSlice: nil,
				errExpected:  true,
			},
		},
		{
			batchSize: 1,
			res: result{
				batchedSlice: nil,
				errExpected:  false,
			},
		},
	}

	t.Log("Testing SplitSnippetSlice function...")
	for testIdx, toTest := range testSet {
		if res, err := SplitSnippetSlice(toTest.slice, toTest.batchSize); err == nil {
			if toTest.res.errExpected {
				t.Errorf("Test <%d> not passed: error expected!", testIdx)
			} else if !compareSnippetSliceBatched(&toTest.res.batchedSlice, &res) {
				t.Errorf("Test <%d> not passed: unexpected result!", testIdx)
			}
		} else if !toTest.res.errExpected {
			t.Errorf("Test <%d> not passed: unexpected error!", testIdx)
		}
	}
}

// =====================================================================

// =====================================================================

// TestFilterSnippetSlice
func TestFilterSnippetSlice(t *testing.T) {
	testSet := []struct {
		snippetSlice []models.Snippet
		filter       []models.Snippet // в функции можно использовать как array, так и slice
		res          []models.Snippet
	}{
		{
			snippetSlice: []models.Snippet{snippet_0, snippet_1, snippet_n5, snippet_3},
			filter:       []models.Snippet{snippet_0, snippet_3},
			res:          []models.Snippet{snippet_1, snippet_n5},
		},
		{
			snippetSlice: []models.Snippet{snippet_0, snippet_1, snippet_0, snippet_3},
			filter:       []models.Snippet{snippet_0, snippet_3},
			res:          []models.Snippet{snippet_1},
		},
		{
			snippetSlice: []models.Snippet{snippet_0, snippet_1, snippet_0, snippet_3},
			filter:       []models.Snippet{},
			res:          []models.Snippet{snippet_0, snippet_1, snippet_0, snippet_3},
		},
		{
			snippetSlice: []models.Snippet{},
			filter:       []models.Snippet{snippet_0, snippet_3},
			res:          []models.Snippet{},
		},
		{
			snippetSlice: []models.Snippet{},
			filter:       []models.Snippet{},
			res:          []models.Snippet{},
		},
		{
			snippetSlice: []models.Snippet{snippet_empty},
			filter:       []models.Snippet{},
			res:          []models.Snippet{snippet_empty},
		},
		{
			snippetSlice: []models.Snippet{snippet_empty},
			filter:       []models.Snippet{snippet_1},
			res:          []models.Snippet{snippet_empty},
		},
		{
			snippetSlice: []models.Snippet{snippet_empty, snippet_19},
			filter:       []models.Snippet{snippet_1, snippet_empty},
			res:          []models.Snippet{snippet_19},
		},
		{
			snippetSlice: []models.Snippet{snippet_empty},
			filter:       []models.Snippet{snippet_empty},
			res:          []models.Snippet{},
		},
		{
			snippetSlice: []models.Snippet{snippet_1, snippet_8},
			filter:       []models.Snippet{snippet_1, snippet_8},
			res:          []models.Snippet{},
		},
		{
			snippetSlice: []models.Snippet{snippet_1, snippet_8},
			res:          []models.Snippet{snippet_1, snippet_8},
		},
		{
			filter: []models.Snippet{snippet_1, snippet_8},
			res:    nil,
		},
	}

	t.Log("Testing FilterSnippetSlice function...")
	for testIdx, toTest := range testSet {
		if res := FilterSnippetSlice(toTest.snippetSlice, toTest.filter...); !compareSnippetSlices(&toTest.res, &res) {
			t.Errorf("Test <%d> not passed: unexpected result!", testIdx)
		}
	}
}

// =====================================================================

// =====================================================================

// TestSliceToMap
func TestSliceToMap(t *testing.T) {
	type result struct {
		snippetMap  map[uint64]models.Snippet
		errExpected bool
	}

	testSet := []struct {
		slice []models.Snippet
		res   result
	}{
		{
			slice: []models.Snippet{snippet_1, snippet_2, snippet_3, snippet_4, snippet_5, snippet_6, snippet_7, snippet_8},
			res: result{
				snippetMap: map[uint64]models.Snippet{
					snippet_1.Id: snippet_1,
					snippet_2.Id: snippet_2,
					snippet_3.Id: snippet_3,
					snippet_4.Id: snippet_4,
					snippet_5.Id: snippet_5,
					snippet_6.Id: snippet_6,
					snippet_7.Id: snippet_7,
					snippet_8.Id: snippet_8,
				},
				errExpected: false,
			},
		},
		{
			slice: []models.Snippet{snippet_1, snippet_5, snippet_6, snippet_7, snippet_8, snippet_1},
			res: result{
				snippetMap:  map[uint64]models.Snippet{},
				errExpected: true,
			},
		},
		{
			slice: []models.Snippet{},
			res: result{
				snippetMap:  map[uint64]models.Snippet{},
				errExpected: false,
			},
		},
		{
			slice: []models.Snippet{snippet_empty},
			res: result{
				snippetMap: map[uint64]models.Snippet{
					snippet_empty.Id: snippet_empty,
				},
				errExpected: false,
			},
		},
		{
			res: result{
				snippetMap:  nil,
				errExpected: false,
			},
		},
		{},
	}

	t.Log("Testing TestSliceToMap function...")
	for testIdx, toTest := range testSet {
		if res, err := SliceToMap(toTest.slice); err == nil {
			if toTest.res.errExpected {
				t.Errorf("Test <%d> not passed: error expected!", testIdx)
			} else if !compareSnippetMaps(&toTest.res.snippetMap, &res) {
				t.Errorf("Test <%d> not passed: unexpected result!", testIdx)
			}
		} else if !toTest.res.errExpected {
			t.Errorf("Test <%d> not passed: unexpected error!", testIdx)
		}
	}
}

// =====================================================================
