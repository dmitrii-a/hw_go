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
	size int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Next:  l.head,
	}
	l.size++
	if l.head == nil {
		l.head = item
		l.tail = item
		return item
	}
	l.head.Prev = item
	l.head = item
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
		Prev:  l.tail,
	}
	l.size++
	if l.tail == nil {
		l.head = item
		l.tail = item
		return item
	}
	l.tail.Next = item
	l.tail = item
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Next = nil
	i.Prev = nil
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	currHead := l.head
	if i == currHead {
		return
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}
	i.Next = currHead
	i.Prev = nil
	currHead.Prev = i
	l.head = i
}

func NewList() List {
	return new(list)
}
