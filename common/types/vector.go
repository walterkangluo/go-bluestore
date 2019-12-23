package types

//T is meant to signify that an interface is supposed to act like a psuedo generic.
// Mainly mean to undo confusion in the codebase
type T interface{}

//AltT is the same as T, but it is used so funcitons can return alternating types
type AltT interface{}

//Vector is a list based on C++ Vector with functional programming functions
type Vector struct {
	data     []T
	size     int
	preAlloc bool
}

func (v *Vector) CreateVector(_preAlloc bool, _size int) {
	v.data = make([]T, 0)
	v.size = 0
	v.preAlloc = false
}

//Init : returns a copy of a initialized vector.
// _Notes:
//		you should set the vecotr equal to the result of it initializing since it returns a copy
//  	EX. "v = v.Init()" is correct usage of this
// _Return:
//		a initialized vector.
func (v *Vector) Init() {
	v.data = make([]T, 0)
	v.size = 0
	v.preAlloc = false
}

//InitWithSlice : returns a copy of a initialized vector that turns the slice into a vector
//_Parameters:
//		elems: T "The slice you want to convert into a vector"
// _Return:
//		a initialized vector.
func (v *Vector) InitWithSlice(elems []T) {
	v.Init()
	for index := 0; index < len(elems); index++ {
		v.PushBack(elems[index])
	}
}

// only be called without any data
func (v *Vector) Reserve(reserveSize int) {
	if v.size > 0 {
		panic("nu nil data")
	}

	v.data = make([]T, 0, reserveSize)
	v.preAlloc = true
}

// only be called without any data
func (v *Vector) ReSize(reserveSize int) {
	if v.size > reserveSize {
		panic("will discard data")
	}

	if v.size > 0 {
		len := v.size
		copy := v.data
		v.data = make([]T, len, reserveSize)

		for i := 0; i < len; i++ {
			v.data[i] = copy[i]
		}
	}
}

//At : Returns the element of the Vector at a specified index
//_Parameters :
//		index : uint64 "the index of the value you want"
//_Return:
//		the element of the vector at a specified index
func (v *Vector) At(index int) T {
	if index >= v.size || index > cap(v.data) {
		return nil
	}

	return v.data[index]
}

//SetAt : Set the element of the Vector at a specified index
//_Parameters :
//		index : uint64 "the index of the value you want"
//		value : int "The element you want to set the value at a specified index instead"
func (v *Vector) SetAt(index int, value T) {
	if index > cap(v.data) && v.preAlloc {
		panic("out of range")
	}
	v.data[index] = value
}

//Size : Returns the size of the vector
//_Return:
//		the size of the vector
func (v *Vector) Size() int {
	return v.size
}

//PushBack : Adds an item to the end of the vector, just like C++
// _Parameters:
//			element: T "the value you want to add"
func (v *Vector) PushBack(element T) {
	v.data = append(v.data, element)
	v.size++
}

//PopBack : removes an item to the end of the vector, just like C++
func (v *Vector) PopBack() {
	lastIndex := v.Size() - 1
	v.RemoveAt(lastIndex)
}

//RemoveAt : removes an item from the vector at a specified index
// _Parameters:
//			index: int "the index of the value you want to remove it from"
func (v *Vector) RemoveAt(index int) {
	v.data = append(v.data[:index], v.data[index+1:]...)
	v.size--
}

//InsertAt : insert an item from the vector at a specified index
// _Parameters:
//			index: int "the index of the value you want to remove it from"
func (v *Vector) InsertAt(index int, element T) {
	v.data = append(v.data[:index], element, v.data[index:])
	v.size++
}

func (v *Vector) Front() T {
	return v.data[0]
}

func (v *Vector) End() T {
	return v.data[v.size]
}

func (v *Vector) Back() T {
	if v.size == 0 {
		return nil
	}
	return v.data[v.size-1]
}

func (v *Vector) Empty() bool {
	return v.size == 0
}

func (v *Vector) Capacity() int {
	return cap(v.data)
}

func (v *Vector) Assign(element T, index int) {
	v.data[index] = element
}

//Clear : removes all elements from the vector
func (v *Vector) Clear() {
	for v.Size() > 0 {
		v.RemoveAt(0)
	}
}

//SortStruct : sorts the vector
//_Notes: Use this to sort your structs. You must implment a compare funciton that runs like this
// 			func(v1, v2){ return v1.compareValue < v2,compareValue }
// 			You can define an inline function that returns the opposite of your sort in you Struct if you want it to be reversed
// _Parameters:
//			sortOperator: T "the reference to the function that would be used as a comparison operator"
func (v *Vector) SortStruct(function func(v1 T, v2 T) bool) {
	//Golang is picky about uint64 and ints
	for i := 0; i < int(v.size-1); i++ {
		for j := i + 1; j < int(v.size); j++ {
			if function(v.data[i], v.data[j]) && i != j {
				temp := v.data[j]
				v.data[j] = v.data[i]
				v.data[i] = temp
			}
		}
	}
}

//FpMap : Returns a copy of the vector that can be altered or not through the inline function
//	_Parameters:
//		function: T "A funciton with one arugument that matches the datatype of the vector's elements,
// it must return a value that is the same datatype as the slice values or it will crash the program"
//	_Returns:
//		T "The filtered vector that is the same datatype as the vector"
//
func (v *Vector) FpMap(function func(T) T) Vector {
	var dataCopy Vector

	dataCopy.Init()

	for index := 0; index < len(v.data); index++ {
		if function != nil {
			mapFunc := function(v.data[index]) //v.Mapable(v.data[index])
			dataCopy.PushBack(mapFunc)
		} else {
			dataCopy.PushBack(v.data[index])
		}
	}
	return dataCopy
}

//FpReduce : reduces a whole value into one and treat it as a single sum from the vector
//	_Parameters:
//		function: T "A funciton with one arugument that matches the datatype of the vector's elements,
// it must return a value that is the same datatype as the slice values or it will crash the program
//	_Returns:
//		T "The accumulate value that is the same datatype as the vector"
//
func (v *Vector) FpReduce(function func(T, T) T) T {
	sum := v.data[0]
	for index := 1; index < len(v.data); index++ {
		if function != nil {
			sum = function(sum, v.data[index])
		} else {
			switch sum.(type) {
			//For those of you wondering why there is a sum temp, it is because
			//Go doesn't allow use to += with sum.(int), which is something we are trying to do
			//THis seems like a little hack, but at least not a painful one....
			case int:
				sumTemp := sum.(int)
				sumTemp += v.data[index].(int)
				sum = sumTemp
				break
			case int16:
				sumTemp := sum.(int16)
				sumTemp += v.data[index].(int16)
				sum = sumTemp
			//The documentation said the int32 is not an alias for int
			case int32:
				sumTemp := sum.(int32)
				sumTemp += v.data[index].(int32)
				sum = sumTemp
			case int64:
				sumTemp := sum.(int64)
				sumTemp += v.data[index].(int64)
				sum = sumTemp
			case float32:
				sumTemp := sum.(float32)
				sumTemp += v.data[index].(float32)
				sum = sumTemp
			case float64:
				sumTemp := sum.(float64)
				sumTemp += v.data[index].(float64)
				sum = sumTemp
			case string:
				sumTemp := sum.(string)
				sumTemp += v.data[index].(string)
				sum = sumTemp
			default:
				panic(`Nil function failed to accumulate vector values, possible reasons \n
				* Vector has mismatching types \n
				* Attempting to add Types that don't have the + operator supported by default in Golang`)
			}
		}
	}
	return sum
}

//FpFilter : Returns a filtered copy of the vector that satifys the given conditions
//	_Parameters:
//		function: T "A funciton with one arugument that matches the datatype of the vector's elements,
// it must return a boolean or it will crash the program"
//	_Returns:
//		T "The filtered vector that is the same datatype as the vector"
//
func (v *Vector) FpFilter(function func(T) bool) Vector {
	var dataCopy Vector
	dataCopy.Init()
	for index := 0; index < len(v.data); index++ {
		if function(v.data[index]) {
			dataCopy.PushBack(v.data[index])
		}
	}
	return dataCopy
}

//FpIndexOf : returns the index of first occurance of the element if found, ortherwise it returns the size of the vecotor
// _Parameters:
//		elementToFind: T "The element that is desired to be found"
// _Return :
//		The index of the desired element
func (v *Vector) FpIndexOf(function func(T) bool) int {
	for index := 0; index < v.Size()-1; index++ {
		if function(v.data[index]) {
			return index
		}
	}
	// This is like saying "last" from the normal std::find
	return v.Size()
}
