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
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}

	if l.front != nil {
		l.front.Prev = newItem
	} else {
		l.back = newItem
	}

	l.front = newItem
	l.length++
	return newItem
}

func (l *list) PushBack(v any) *ListItem {
	newItem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.back != nil {
		l.back.Next = newItem
	} else {
		l.front = newItem
	}

	l.back = newItem
	l.length++
	return newItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.front = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.back = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}

	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return &list{}
}
