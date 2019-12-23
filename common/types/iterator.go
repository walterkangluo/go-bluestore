package types

import "container/list"

type Iterator interface {
	HasNext() bool
	Value() interface{}
	Next()
}

type ListIterator struct {
	cur *list.Element
	end *list.Element
}

func (this *ListIterator) HasNext() bool {
	return this.cur != this.end
}

func (this *ListIterator) Next() {
	this.cur = this.cur.Next()
}

func (this *ListIterator) Value() interface{} {
	return this.cur.Value
}

type Container interface {
	Iterator() Iterator
}

type List struct {
	list list.List
}

func (this *List) Iterator() Iterator {
	return &ListIterator{this.list.Front(), this.list.Back()}
}

func (this *List) Add(value interface{}) {
	this.list.PushBack(value)
}
