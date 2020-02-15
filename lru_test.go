package lru

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSet(t *testing.T) {
	size := 3
	lru := NewLRU(size)

	for i := 0; i < size; i++ {
		lru.Set(strconv.Itoa(i), i)
	}

	e, isExist := lru.Get("1")
	if !isExist {
		t.Fatalf("failed to get value")
	}

	if got := e.(int); got != 1 {
		t.Errorf("value got=%d; want=%d", got, 1)
	}

	lru.Set("3", 3)

	_, isExist = lru.Get("0")
	if isExist {
		t.Fatalf("the oldest element is presented")
	}

	ok := lru.Remove("1")
	if !ok {
		t.Fatalf("failed to remove element")
	}

	if _, isExist := lru.Get("1"); isExist {
		t.Fatalf("the element wasn't removed")
	}

	fmt.Printf("%+v\n", lru)

}
