package types

type S map[interface{}]interface{}

type Set struct {
	data []S
	size uint32
}

/*
begin()     　　 ,返回set容器的第一个元素

end() 　　　　 ,返回set容器的最后一个元素

clear()   　　     ,删除set容器中的所有的元素

empty() 　　　,判断set容器是否为空

max_size() 　 ,返回set容器可能包含的元素最大个数

size() 　　　　 ,返回当前set容器中的元素个数

rbegin　　　　 ,返回的值和end()相同

rend()　　　　 ,返回的值和rbegin()相同
*/

func (s *Set) Init() {
	s.data = make([]S, 0)
	s.size = 0
}

func (s *Set) Size() uint32 {
	return s.size
}

func (s *Set) Begin() S {
	if 0 == s.Size() {
		return nil
	}
	return s.data[0]
}

func (s *Set) End() S {
	if 0 == s.Size() {
		return nil
	}
	return s.data[s.size]
}

func (s *Set) Empty() bool {

	return s.Size() == 0
}

func (s *Set) Back() S {
	if 0 == s.Size() {
		return nil
	}
	return s.data[s.size-1]
}

func (s *Set) Push(ele S) {
	s.data = append(s.data, ele)
	s.size++
}

func (s *Set) PoP(ele S) S {
	if s.Empty() {
		return nil
	}
	i := s.Size()
	s.size--
	return s.data[i]
}
