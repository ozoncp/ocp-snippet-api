package utils

import (
	"testing"

	models "github.com/ozoncp/ocp-snippet-api/internal/snippet"
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
	snippet_11    = models.Snippet{Id: 11}
	snippet_19    = models.Snippet{Id: 19}
	snippet_0     = models.Snippet{Id: 0}
	snippet_empty = models.Snippet{}
)

// =====================================================================

// =====================================================================
// Compare functions:

func compareSnippetSlices(l *models.Snippets, r *models.Snippets) bool {
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

// compareSnippetSliceBatched - функция сравнения слайсов батчей ([]models.Snippets).
// Считает слайсы батчей равными, если их длины равны, длины батчей равны и указатели на сниппеты либо значения снипеттов равны
// (тестируемая функция не делает глубокого копирования, однако функция сравнения рассматривает и такой кейс).
func compareSnippetSliceBatched(l *[]models.Snippets, r *[]models.Snippets) bool {
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
func compareReversedSnippetMaps(l *map[*models.Snippet]uint64, r *map[*models.Snippet]uint64) bool {
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
func compareSnippetMaps(l *models.SnippetMap, r *models.SnippetMap) bool {
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
		batchedSlice []models.Snippets
		errExpected  bool
	}

	testSet := []struct {
		slice     models.Snippets
		batchSize uint
		res       result
	}{
		{
			slice:     models.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 3,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1, &snippet_2, &snippet_3}, {&snippet_4, &snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1},
			batchSize: 3,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1, &snippet_2}, {&snippet_3, &snippet_4}, {&snippet_5, &snippet_6}, {&snippet_7, &snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8, &snippet_9},
			batchSize: 5,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5}, {&snippet_6, &snippet_7, &snippet_8, &snippet_9}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{},
			batchSize: 5,
			res: result{
				batchedSlice: []models.Snippets{{}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1, &snippet_3, &snippet_8},
			batchSize: 2,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1, &snippet_3}, {&snippet_8}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_1, &snippet_3, &snippet_8, &snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_empty},
			batchSize: 5,
			res: result{
				batchedSlice: []models.Snippets{{&snippet_empty}},
				errExpected:  false,
			},
		},
		{
			slice:     models.Snippets{&snippet_1, &snippet_3, &snippet_8, &snippet_empty},
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
		snippetMap    models.SnippetMap
		res           map[*models.Snippet]uint64
		panicExpected bool
	}{
		{
			snippetMap: models.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_8,
				4: &snippet_n5,
				5: &snippet_11,
				6: &snippet_19,
			},
			res: map[*models.Snippet]uint64{
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
			snippetMap: models.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_n5,
			},
			res: map[*models.Snippet]uint64{
				&snippet_1:  1,
				&snippet_3:  2,
				&snippet_n5: 3,
			},
			panicExpected: false,
		},
		{
			snippetMap: models.SnippetMap{
				1: &snippet_1,
				2: &snippet_3,
				3: &snippet_1,
			},
			res:           map[*models.Snippet]uint64{},
			panicExpected: true,
		},
		{
			snippetMap:    models.SnippetMap{},
			res:           map[*models.Snippet]uint64{},
			panicExpected: false,
		},
		{
			snippetMap: models.SnippetMap{
				1: &snippet_1,
			},
			res: map[*models.Snippet]uint64{
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
		snippetSlice models.Snippets
		filter       models.Snippets // в функции можно использовать как array, так и slice
		res          models.Snippets
	}{
		{
			snippetSlice: models.Snippets{&snippet_0, &snippet_1, &snippet_n5, &snippet_3},
			filter:       models.Snippets{&snippet_0, &snippet_3},
			res:          models.Snippets{&snippet_1, &snippet_n5},
		},
		{
			snippetSlice: models.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       models.Snippets{&snippet_0, &snippet_3},
			res:          models.Snippets{&snippet_1},
		},
		{
			snippetSlice: models.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
			filter:       models.Snippets{},
			res:          models.Snippets{&snippet_0, &snippet_1, &snippet_0, &snippet_3},
		},
		{
			snippetSlice: models.Snippets{},
			filter:       models.Snippets{&snippet_0, &snippet_3},
			res:          models.Snippets{},
		},
		{
			snippetSlice: models.Snippets{},
			filter:       models.Snippets{},
			res:          models.Snippets{},
		},
		{
			snippetSlice: models.Snippets{&snippet_empty},
			filter:       models.Snippets{},
			res:          models.Snippets{&snippet_empty},
		},
		{
			snippetSlice: models.Snippets{&snippet_empty},
			filter:       models.Snippets{&snippet_1},
			res:          models.Snippets{&snippet_empty},
		},
		{
			snippetSlice: models.Snippets{&snippet_empty, &snippet_19},
			filter:       models.Snippets{&snippet_1, &snippet_empty},
			res:          models.Snippets{&snippet_19},
		},
		{
			snippetSlice: models.Snippets{&snippet_empty},
			filter:       models.Snippets{&snippet_empty},
			res:          models.Snippets{},
		},
		{
			snippetSlice: models.Snippets{&snippet_1, &snippet_8},
			filter:       models.Snippets{&snippet_1, &snippet_8},
			res:          models.Snippets{},
		},
		{
			snippetSlice: models.Snippets{&snippet_1, &snippet_8},
			res:          models.Snippets{&snippet_1, &snippet_8},
		},
		{
			filter: models.Snippets{&snippet_1, &snippet_8},
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
		snippetMap  models.SnippetMap
		errExpected bool
	}

	testSet := []struct {
		slice models.Snippets
		res   result
	}{
		{
			slice: models.Snippets{&snippet_1, &snippet_2, &snippet_3, &snippet_4, &snippet_5, &snippet_6, &snippet_7, &snippet_8},
			res: result{
				snippetMap: models.SnippetMap{
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
			slice: models.Snippets{&snippet_1, &snippet_5, &snippet_6, &snippet_7, &snippet_8, &snippet_1},
			res: result{
				snippetMap:  models.SnippetMap{},
				errExpected: true,
			},
		},
		{
			slice: models.Snippets{},
			res: result{
				snippetMap:  models.SnippetMap{},
				errExpected: false,
			},
		},
		{
			slice: models.Snippets{&snippet_empty},
			res: result{
				snippetMap: models.SnippetMap{
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
