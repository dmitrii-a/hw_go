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

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())
		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("correct order with PushFront", func(t *testing.T) {
		l := NewList()
		l.PushFront(1)
		l.PushFront(2)
		l.PushFront(3)
		require.Equal(t, 3, l.Len())
		item := l.Front()
		require.Equal(t, 3, item.Value)
		require.Equal(t, 2, item.Next.Value)
		require.Equal(t, 1, item.Next.Next.Value)
	})

	t.Run("correct order with PushBack", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		l.PushBack(2)
		l.PushBack(3)
		require.Equal(t, 3, l.Len())
		item := l.Front()
		require.Equal(t, 1, item.Value)
		require.Equal(t, 2, item.Next.Value)
		require.Equal(t, 3, item.Next.Next.Value)
	})

	t.Run("correct order with PushBack/PushFront", func(t *testing.T) {
		l := NewList()
		l.PushBack(1)
		l.PushFront(2)
		l.PushBack(3)
		require.Equal(t, 3, l.Len())
		item := l.Front()
		require.Equal(t, 2, item.Value)
		require.Equal(t, 1, item.Next.Value)
		require.Equal(t, 3, item.Next.Next.Value)
	})

	t.Run("check MoveToFront", func(t *testing.T) {
		l := NewList()
		item := l.PushFront(1)
		l.PushFront(2)
		l.PushFront(3)
		require.Equal(t, 3, l.Front().Value)
		l.MoveToFront(item)
		require.Equal(t, 1, l.Front().Value)
	})

	t.Run("check remove", func(t *testing.T) {
		l := NewList()
		item := l.PushFront(1)
		item2 := l.PushFront(2)
		item3 := l.PushFront(3)
		require.Equal(t, 3, l.Len())
		l.Remove(item2)
		require.Equal(t, 2, l.Len())
		require.Equal(t, (*ListItem)(nil), item2.Next)
		require.Equal(t, (*ListItem)(nil), item2.Prev)
		require.Equal(t, (*ListItem)(nil), item.Next)
		require.Equal(t, (*ListItem)(nil), item3.Prev)
	})

	t.Run("check changing head/tail when remove item", func(t *testing.T) {
		l := NewList()
		item := l.PushFront(1)
		item2 := l.PushFront(2)
		require.Equal(t, 2, l.Len())
		l.Remove(item)
		l.Remove(item2)
		require.Equal(t, 0, l.Len())
		item3 := l.PushFront(3)
		require.Equal(t, (*ListItem)(nil), item3.Next)
		l.Remove(item3)
		item4 := l.PushBack(4)
		require.Equal(t, (*ListItem)(nil), item4.Prev)
	})
}
