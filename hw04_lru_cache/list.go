package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}
func (l *list) Front() *ListItem {
	return l.head
}
func (l *list) Back() *ListItem {
	return l.tail
}
func (l *list) PushFront(v interface{}) *ListItem {
	l.len++
	a := new(ListItem)
	a.Value = v
	a.Next = l.head
	if l.head != nil {
		l.head.Prev = a
	}
	if l.tail == nil {
		l.tail = a
	}

	l.head = a
	return a
}
func (l *list) PushBack(v interface{}) *ListItem {
	l.len++
	a := new(ListItem)
	a.Value = v
	a.Prev = l.tail
	if l.tail != nil {
		l.tail.Next = a
	}
	if l.head == nil {
		l.head = a
	}
	l.tail = a
	return a
}
func (l *list) Remove(i *ListItem) {
	l.len--
	if i.Prev == nil && i.Next == nil {
		l.head = nil
		l.tail = nil
	} else if i.Prev == nil {
		l.head = i.Next
		l.head.Prev = nil
	} else if i.Next == nil {
		l.tail = i.Prev
		l.tail.Next = nil
	} else {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	if i.Next == nil {
		l.tail = i.Prev
		l.tail.Next = nil
	} else {
		i.Prev.Next, i.Next.Prev = i.Next, i.Prev
	}
	i.Next = l.head
	l.head.Prev = i
	i.Prev = nil
	l.head = i

}

func NewList() List {
	return new(list)
}
