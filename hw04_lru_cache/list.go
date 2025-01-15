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
	first  *ListItem
	last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	first := l.first
	newItem := ListItem{v, first, nil}
	if first != nil {
		first.Prev = &newItem
	}
	l.first = &newItem
	l.length++
	if l.length == 1 {
		l.last = &newItem
	}
	return l.first
}

func (l *list) PushBack(v interface{}) *ListItem {
	last := l.last
	newItem := ListItem{v, nil, last}
	if last != nil {
		last.Next = &newItem
	}
	l.last = &newItem
	l.length++
	if l.length == 1 {
		l.first = &newItem
	}
	return l.last
}

func (l *list) Remove(i *ListItem) {
	prev := i.Prev
	next := i.Next
	if prev != nil {
		prev.Next = next
	} else {
		l.first = next
	}
	if next != nil {
		next.Prev = prev
	} else {
		l.last = prev
	}
	if l.length > 0 {
		l.length--
	}
	i.Next = nil
	i.Prev = nil
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.length++
	if l.first != nil {
		l.first.Prev = i
	}
	i.Next = l.first
	l.first = i
}

func NewList() List {
	return new(list)
}
