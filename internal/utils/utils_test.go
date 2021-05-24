package utils

import (
	"testing"
)

// =====================================================================
// Testing snippets:
var (
	snippet_1     = Snippet{Field: 1}
	snippet_2     = Snippet{Field: 2}
	snippet_3     = Snippet{Field: 3}
	snippet_4     = Snippet{Field: 4}
	snippet_5     = Snippet{Field: 5}
	snippet_6     = Snippet{Field: 6}
	snippet_7     = Snippet{Field: 7}
	snippet_8     = Snippet{Field: 8}
	snippet_9     = Snippet{Field: 9}
	snippet_n5    = Snippet{Field: -5}
	snippet_11    = Snippet{Field: 11}
	snippet_19    = Snippet{Field: 19}
	snippet_0     = Snippet{Field: 0}
	snippet_empty = Snippet{}
)

// =====================================================================

// =====================================================================
// Compare functions

// compareSnippets
// Сравнивает указатели на Snippet и сравнивает по значению (сравнение по значению д.б. в реализации Snippet)
func compareSnippets(l *Snippet, r *Snippet) bool {
	return l == r || l.Field == r.Field
}
func compareSnippetSlices(l *SnippetSlice, r *SnippetSlice) bool {
	if len(*l) != len(*r) {
		return false
	}

	for idx, lSnippet := range *l {
		if !compareSnippets(lSnippet, (*r)[idx]) {
			return false
		}
	}

	return true
}

// compareSnippetSliceBatched - функция сравнения слайсов батчей (SnippetSliceBatched).
// Считает слайсы батчей равными, если их длины равны, длины батчей равны и указатели на сниппеты либо значения снипеттов равны
// (тестируемая функция не делает глубокого копирования, однако функция сравнения рассматривает и такой кейс).
func compareSnippetSliceBatched(l *SnippetSliceBatched, r *SnippetSliceBatched) bool {
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
func compareReversedSnippetMaps(l *ReversedSnippetMap, r *ReversedSnippetMap) bool {
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

// =====================================================================

// =====================================================================
// TestSplitSnippetSlice

func TestSplitSnippetSlice(t *testing.T) {
	type result struct {
		batchedSlice SnippetSliceBatched
		errExpected  bool
	}

	testSet := []struct {
		slice     SnippetSlice
		batchSize int
		res       result
	}{
		{
			slice:     SnippetSlice{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 3,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1, &snippet_2, &snippet_3}, {&snippet_4, &snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1},
			batchSize: 3,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1, &snippet_2}, {&snippet_3, &snippet_4}, {&snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8, &snippet_9},
			batchSize: 5,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5}, {&snippet_6, &snippet_7, &snippet_8, &snippet_9}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{},
			batchSize: 5,
			res: result{
				batchedSlice: SnippetSliceBatched{},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1, &snippet_3, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1, &snippet_3}, {&snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_1, &snippet_3, &snippet_8, &snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: SnippetSliceBatched{{&snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     SnippetSlice{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
			batchSize: 0,
			res: result{
				batchedSlice: SnippetSliceBatched{},
				errExpected:  true,
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
		snippetMap    SnippetMap
		res           ReversedSnippetMap
		panicExpected bool
	}{
		{
			snippetMap: SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_8,
				4: &snippet_n5,
				5: &snippet_11,
				6: &snippet_19,
			},
			res: ReversedSnippetMap{
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
			snippetMap: SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_n5,
			},
			res: ReversedSnippetMap{
				&snippet_1:  1,
				&snippet_3:  2,
				&snippet_n5: 3,
			},
			panicExpected: false,
		},
		{
			snippetMap: SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_1,
			},
			res:           ReversedSnippetMap{},
			panicExpected: true,
		},
		{
			snippetMap:    SnippetMap{},
			res:           ReversedSnippetMap{},
			panicExpected: false,
		},
		{
			snippetMap: SnippetMap{
				1: &snippet_1,
			},
			res: ReversedSnippetMap{
				&snippet_1: 1,
			},
			panicExpected: false,
		},
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
		snippetSlice SnippetSlice
		filter       SnippetSlice // в функции можно использовать как array, так и slice
		res          SnippetSlice
	}{
		{
			snippetSlice: SnippetSlice{&snippet_0, &snippet_1, &snippet_n5, &snippet_3},
			filter:       SnippetSlice{&snippet_0, &snippet_3},
			res:          SnippetSlice{&snippet_1, &snippet_n5},
		},
		{
			snippetSlice: SnippetSlice{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       SnippetSlice{&snippet_0, &snippet_3},
			res:          SnippetSlice{&snippet_1},
		},
		{
			snippetSlice: SnippetSlice{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       SnippetSlice{},
			res:          SnippetSlice{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
		},
		{
			snippetSlice: SnippetSlice{},
			filter:       SnippetSlice{&snippet_0, &snippet_3},
			res:          SnippetSlice{},
		},
		{
			snippetSlice: SnippetSlice{},
			filter:       SnippetSlice{},
			res:          SnippetSlice{},
		},
		{
			snippetSlice: SnippetSlice{&snippet_empty},
			filter:       SnippetSlice{},
			res:          SnippetSlice{&snippet_empty},
		},
		{
			snippetSlice: SnippetSlice{&snippet_empty},
			filter:       SnippetSlice{&snippet_1},
			res:          SnippetSlice{&snippet_empty},
		},
		{
			snippetSlice: SnippetSlice{&snippet_empty, &snippet_19},
			filter:       SnippetSlice{&snippet_1, &snippet_empty},
			res:          SnippetSlice{&snippet_19},
		},
		{
			snippetSlice: SnippetSlice{&snippet_empty},
			filter:       SnippetSlice{&snippet_empty},
			res:          SnippetSlice{},
		},
		{
			snippetSlice: SnippetSlice{&snippet_1, &snippet_8},
			filter:       SnippetSlice{&snippet_1, &snippet_8},
			res:          SnippetSlice{},
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
