package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()
		// [20] - > [10, 20]
		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		l.Remove(l.Front())
		l.Remove(l.Front())
		require.Equal(t, 0, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 5, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)
		a := l.Back()
		l.MoveToFront(a)
		l.MoveToFront(a) // [80, 60, 40, 10, 30, 50, 70]
		// [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, l.Len())
		elems2 := make([]int, l.Len())
		for j, i := 0, l.Front(); i != nil; j, i = j+1, i.Next {
			elems[j] = i.Value.(int)
		}
		for j, i := l.Len()-1, l.Back(); i != nil; i, j = i.Prev, j-1 {
			elems2[j] = i.Value.(int)
		}

		require.Equal(t, elems2, elems)
		require.Equal(t, []int{70, 80, 60, 40, 50}, elems)
	})
}
