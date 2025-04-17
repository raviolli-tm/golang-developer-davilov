package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func check(t *testing.T, testList List, expected []int) {
	t.Helper()
	i := testList.Front()
	for j := 0; j < testList.Len(); j++ {
		i = i.Next
	}
	require.Nil(t, i)

	i = testList.Back()
	for j := 0; j < testList.Len(); j++ {
		i = i.Prev
	}
	require.Nil(t, i)

	elemsForward := make([]int, testList.Len())
	elemsBackward := make([]int, testList.Len())
	for j, i := 0, testList.Front(); i != nil; j, i = j+1, i.Next {
		elemsForward[j] = i.Value.(int)
	}
	for j, i := testList.Len()-1, testList.Back(); i != nil; i, j = i.Prev, j-1 {
		elemsBackward[j] = i.Value.(int)
	}

	require.Equal(t, elemsForward, elemsBackward)
	require.Equal(t, expected, elemsForward)
}

func TestList(t *testing.T) {
	testList := NewList()

	t.Run("empty list", func(t *testing.T) {
		require.Equal(t, 0, testList.Len())
		require.Nil(t, testList.Front())
		require.Nil(t, testList.Back())
	})

	t.Run("complex", func(t *testing.T) {
		testList.PushFront(10)
		testList.PushBack(20)
		check(t, testList, []int{10, 20})
		testList.PushBack(30)

		require.Equal(t, 3, testList.Len())

		check(t, testList, []int{10, 20, 30})
	})

	t.Run("remove", func(t *testing.T) {
		middle := testList.Front().Next
		testList.Remove(middle)
		check(t, testList, []int{10, 30})
		testList.Remove(testList.Front())
		testList.Remove(testList.Front())
		check(t, testList, []int{})
		require.Equal(t, 0, testList.Len())
		require.Nil(t, testList.Back())
		require.Nil(t, testList.Front())
	})

	t.Run("move to front", func(t *testing.T) {
		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				testList.PushFront(v)
			} else {
				testList.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 5, testList.Len())
		require.Equal(t, 80, testList.Front().Value)
		require.Equal(t, 70, testList.Back().Value)
		a := testList.Back()
		testList.MoveToFront(a)
		testList.MoveToFront(a)
		check(t, testList, []int{70, 80, 60, 40, 50})
		require.Equal(t, 5, testList.Len())
		a = testList.Back().Prev
		testList.MoveToFront(a)
		check(t, testList, []int{40, 70, 80, 60, 50})
		require.Equal(t, 5, testList.Len())
	})
}
