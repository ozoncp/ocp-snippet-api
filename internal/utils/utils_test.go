package utils

import (
	"testing"

	"github.com/ozoncp/ocp-snippet-api/internal/snippet"
)

// =====================================================================
// Testing snippets:
var (
	snippet_1     = snippet.Snippet{Id: 1}
	snippet_2     = snippet.Snippet{Id: 2}
	snippet_3     = snippet.Snippet{Id: 3}
	snippet_4     = snippet.Snippet{Id: 4}
	snippet_5     = snippet.Snippet{Id: 5}
	snippet_6     = snippet.Snippet{Id: 6}
	snippet_7     = snippet.Snippet{Id: 7}
	snippet_8     = snippet.Snippet{Id: 8}
	snippet_9     = snippet.Snippet{Id: 9}
	snippet_n5    = snippet.Snippet{Id: 15}
	snippet_11    = snippet.Snippet{Id: 11}
	snippet_19    = snippet.Snippet{Id: 19}
	snippet_0     = snippet.Snippet{Id: 0}
	snippet_empty = snippet.Snippet{}
)

// =====================================================================

// =====================================================================
// Compare functions:

func compareSnippetSlices(l *snippet.Snippets, r *snippet.Snippets) bool {
	if len(*l) != len(*r) {
		return false
	}

	for idx, lSnippet := range *l {
		if !(lSnippet == (*r)[idx] || *lSnippet == *(*r)[idx]) {
			return false
		}
	}

	return true
}

// compareSnippetSliceBatched - функция сравнения слайсов батчей ([]snippet.Snippets).
// Считает слайсы батчей равными, если их длины равны, длины батчей равны и указатели на сниппеты либо значения снипеттов равны
// (тестируемая функция не делает глубокого копирования, однако функция сравнения рассматривает и такой кейс).
func compareSnippetSliceBatched(l *[]snippet.Snippets, r *[]snippet.Snippets) bool {
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
func compareReversedSnippetMaps(l *map[*snippet.Snippet]uint64, r *map[*snippet.Snippet]uint64) bool {
	if len(*l) != len(*r) {
		return false
	}

	for lKey, lValue := range *l {
		if rValue, found := (*r)[lKey]; !found || rValue != lValue {
			return false
		}
	}

	return true
}
func compareSnippetMaps(l *snippet.SnippetMap, r *snippet.SnippetMap) bool {
	if len(*l) != len(*r) {
		return false
	}

	for lKey, lValue := range *l {
		if rValue, found := (*r)[lKey]; !found || !(lValue == rValue || *lValue == *rValue) {
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
		batchedSlice []snippet.Snippets
		errExpected  bool
	}

	testSet := []struct {
		slice     snippet.Snippets
		batchSize uint
		res       result
	}{
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 3,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1, &snippet_2, &snippet_3}, {&snippet_4, &snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1},
			batchSize: 3,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1, &snippet_2}, {&snippet_3, &snippet_4}, {&snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8, &snippet_9},
			batchSize: 5,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5}, {&snippet_6, &snippet_7, &snippet_8, &snippet_9}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{},
			batchSize: 5,
			res: result{
				batchedSlice: []snippet.Snippets{{}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_3, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1, &snippet_3}, {&snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_1, &snippet_3, &snippet_8, &snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: []snippet.Snippets{{&snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     snippet.Snippets{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
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
		if res, err := SplitSnippetSlice(&toTest.slice, toTest.batchSize); err == nil {
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

// TestReverseSnippetMap
func TestReverseSnippetMap(t *testing.T) {
	testSet := []struct {
		snippetMap    snippet.SnippetMap
		res           map[*snippet.Snippet]uint64
		panicExpected bool
	}{
		{
			snippetMap: snippet.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_8,
				4: &snippet_n5,
				5: &snippet_11,
				6: &snippet_19,
			},
			res: map[*snippet.Snippet]uint64{
				&snippet_1:  1,
				&snippet_3:  2,
				&snippet_8:  3,
				&snippet_n5: 4,
				&snippet_11: 5,
				&snippet_19: 6,
			},
			panicExpected: false,
		},
		{
			snippetMap: snippet.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_n5,
			},
			res: map[*snippet.Snippet]uint64{
				&snippet_1:  1,
				&snippet_3:  2,
				&snippet_n5: 3,
			},
			panicExpected: false,
		},
		{
			snippetMap: snippet.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_1,
			},
			res:           map[*snippet.Snippet]uint64{},
			panicExpected: true,
		},
		{
			snippetMap:    snippet.SnippetMap{},
			res:           map[*snippet.Snippet]uint64{},
			panicExpected: false,
		},
		{
			snippetMap: snippet.SnippetMap{
				1: &snippet_1,
			},
			res: map[*snippet.Snippet]uint64{
				&snippet_1: 1,
			},
			panicExpected: false,
		},
		{},
	}

	testIdxPtr := new(int) // указатель, т.к. этот индекс используется при обработки паники в функторе, а замыкание происходит в момент инициализации функтора
	panicHandler := func() {
		if testIdxPtr == nil {
			return
		}
		t.Logf("%d", *testIdxPtr)
		if obj := recover(); obj != nil && !testSet[*testIdxPtr].panicExpected {
			t.Errorf("Test <%d> not passed: unexpected panic!", *testIdxPtr)
		}
	}

	t.Log("Testing ReverseSnippetMap function...")

	for testIdx, toTest := range testSet {
		testIdxPtr = &testIdx

		defer panicHandler()

		res := ReverseSnippetMap(&toTest.snippetMap)

		if toTest.panicExpected {
			t.Errorf("Test <%d> not passed: panic expected!", *testIdxPtr)
		} else if !compareReversedSnippetMaps(&toTest.res, &res) {
			t.Errorf("Test <%d> not passed: unexpected result!", *testIdxPtr)
		}
	}
	testIdxPtr = nil
}

// =====================================================================

// =====================================================================

// TestFilterSnippetSlice
func TestFilterSnippetSlice(t *testing.T) {
	testSet := []struct {
		snippetSlice snippet.Snippets
		filter       snippet.Snippets // в функции можно использовать как array, так и slice
		res          snippet.Snippets
	}{
		{
			snippetSlice: snippet.Snippets{&snippet_0, &snippet_1, &snippet_n5, &snippet_3},
			filter:       snippet.Snippets{&snippet_0, &snippet_3},
			res:          snippet.Snippets{&snippet_1, &snippet_n5},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       snippet.Snippets{&snippet_0, &snippet_3},
			res:          snippet.Snippets{&snippet_1},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       snippet.Snippets{},
			res:          snippet.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
		},
		{
			snippetSlice: snippet.Snippets{},
			filter:       snippet.Snippets{&snippet_0, &snippet_3},
			res:          snippet.Snippets{},
		},
		{
			snippetSlice: snippet.Snippets{},
			filter:       snippet.Snippets{},
			res:          snippet.Snippets{},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_empty},
			filter:       snippet.Snippets{},
			res:          snippet.Snippets{&snippet_empty},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_empty},
			filter:       snippet.Snippets{&snippet_1},
			res:          snippet.Snippets{&snippet_empty},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_empty, &snippet_19},
			filter:       snippet.Snippets{&snippet_1, &snippet_empty},
			res:          snippet.Snippets{&snippet_19},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_empty},
			filter:       snippet.Snippets{&snippet_empty},
			res:          snippet.Snippets{},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_1, &snippet_8},
			filter:       snippet.Snippets{&snippet_1, &snippet_8},
			res:          snippet.Snippets{},
		},
		{
			snippetSlice: snippet.Snippets{&snippet_1, &snippet_8},
			res:          snippet.Snippets{&snippet_1, &snippet_8},
		},
		{
			filter: snippet.Snippets{&snippet_1, &snippet_8},
			res:    nil,
		},
	}

	t.Log("Testing FilterSnippetSlice function...")
	for testIdx, toTest := range testSet {
		if res := FilterSnippetSlice(&toTest.snippetSlice, toTest.filter...); !compareSnippetSlices(&toTest.res, &res) {
			t.Errorf("Test <%d> not passed: unexpected result!", testIdx)
		}
	}
}

// =====================================================================

// =====================================================================

// TestSliceToMap
func TestSliceToMap(t *testing.T) {
	type result struct {
		snippetMap  snippet.SnippetMap
		errExpected bool
	}

	testSet := []struct {
		slice snippet.Snippets
		res   result
	}{
		{
			slice: snippet.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			res: result{
				snippetMap: snippet.SnippetMap{
					snippet_1.Id: &snippet_1,
					snippet_2.Id: &snippet_2,
					snippet_3.Id: &snippet_3,
					snippet_4.Id: &snippet_4,
					snippet_5.Id: &snippet_5,
					snippet_6.Id: &snippet_6,
					snippet_7.Id: &snippet_7,
					snippet_8.Id: &snippet_8,
				},
				errExpected: false,
			},
		},
		{
			slice: snippet.Snippets{&snippet_1, &snippet_5, &snippet_6, &snippet_7, &snippet_8, &snippet_1},
			res: result{
				snippetMap:  snippet.SnippetMap{},
				errExpected: true,
			},
		},
		{
			slice: snippet.Snippets{},
			res: result{
				snippetMap:  snippet.SnippetMap{},
				errExpected: false,
			},
		},
		{
			slice: snippet.Snippets{&snippet_empty},
			res: result{
				snippetMap: snippet.SnippetMap{
					snippet_empty.Id: &snippet_empty,
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
		if res, err := SliceToMap(&toTest.slice); err == nil {
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
