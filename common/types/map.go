package types

import (
	"container/list"
)

type Keyer interface {
	GetKey() interface{}
	GetVal() interface{}
	ModifyVal(interface{})
}

type MapList struct {
	dataMap  map[interface{}]*list.Element
	dataList *list.List
}

func NewMapList() *MapList {
	return &MapList{
		dataMap:  make(map[interface{}]*list.Element),
		dataList: list.New(),
	}
}

func (mapList *MapList) Exists(data Keyer) (Keyer, bool) {
	element, exists := mapList.dataMap[data.GetKey()]
	if exists {
		return element.Value.(Keyer), true
	}
	return nil, false
}

func (mapList *MapList) Push(data Keyer) bool {
	if _, ok := mapList.Exists(data); ok {
		// already exists
		return false
	}
	elem := mapList.dataList.PushBack(data)
	mapList.dataMap[data.GetKey()] = elem
	return true
}

func (mapList *MapList) TryPush(data Keyer) {
	el, ok := mapList.Exists(data)
	if ok {
		el.(Keyer).ModifyVal(data.GetVal())
	} else {
		elem := mapList.dataList.PushBack(data)
		mapList.dataMap[data.GetKey()] = elem
	}
}

func (mapList *MapList) Remove(data Keyer) {
	if _, ok := mapList.Exists(data); !ok {
		return
	}
	mapList.dataList.Remove(mapList.dataMap[data.GetKey()])
	delete(mapList.dataMap, data.GetKey())
}

func (mapList *MapList) Size() int {
	return mapList.dataList.Len()
}

func (mapList *MapList) Walk(cb func(data Keyer)) {
	for elem := mapList.dataList.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value.(Keyer))
	}
}

type Elements struct {
	key interface{}
	val interface{}
}

func CreateElements(_key interface{}, _val interface{}) *Elements{
	return &Elements{
		key: _key,
		val: _val,
	}
}

func (e Elements) GetKey() interface{} {
	return e.key
}

func (e Elements) GetVal() interface{} {
	return e.val
}

func (e *Elements) ModifyVal(_val interface{}) {
	e.val = _val
}
