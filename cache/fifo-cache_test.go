package cache

import (
	"strconv"
	"testing"
)

type TestItem struct{
	key string
	val int
}

func getItems(num int)[]*TestItem {

	tests := make([]*TestItem, num)

	for i := 0; i < len(tests); i++{
		tests[i] = &TestItem{
			key: strconv.Itoa(i),
			val: i,
		}
	}
	return tests
}

func Test_SetAndFind(t *testing.T){

	gap := 2
	tests := getItems(10)
	cache := NewFifoCache(len(tests) - gap)

	for _, testcase := range tests{
		cache.Set(testcase.key, testcase.val)
	}
	/// the cache should now be filled with the last n - gap test cases (where n is the number of cases)
	/// the first cases, upto gap, should not be in the cache
	i := 0
	for ; i < gap; i++{
		val, found := cache.Find(tests[i].key)
		if found{
			t.Error("Found items that should not be in the cache at ", i, " val ", val)
		}
	}
	///now check the remaining items can be found with the correct value
	for ; i < len(tests); i++{
		val, found := cache.Find(tests[i].key)
		if !found{
			t.Error("Failed to find ", tests[i].key)
		}
		if val.(int) != tests[i].val{
			t.Error("Mismatch value ", val, tests[i].val)
		}
	}
}

func Test_SameItem(t *testing.T){

	tests := getItems(2)
	cache := NewFifoCache(2)

	//// push item 1 onto the cache multiple times, then put item 2, ensure we can get both back
	for i := 0; i < len(tests); i++{
		cache.Set(tests[0].key, tests[0].val)
	}
	for i := 0; i < len(tests); i++{
		cache.Set(tests[i].key, tests[i].val)
	}
	for _, testcase := range tests{
		val, found := cache.Find(testcase.key)
		if !found{
			t.Error("Failed to find ", testcase.key)
		}
		if val.(int) != testcase.val{
			t.Error("Mimatch value ", val, testcase.val)
		}
	}
}
