package evaluator

import (
	"strings"

	. "github.com/BEN00262/simpleLang/exceptions"
	. "github.com/BEN00262/simpleLang/parser"
)

type DictItem struct {
	key   string
	value interface{}
}

type DictEvalNode struct {
	capacity int64 // sum of the windows ( allocated )
	items    []DictItem
}

func NewDict(size int64) DictEvalNode {
	var capacity int64 = 100

	if size >= 100 {
		capacity *= 2
	}

	return DictEvalNode{
		items:    make([]DictItem, capacity),
		capacity: capacity,
	}
}

func knuthHashing(value int64, size int64) int64 {
	if size >= 0 && size <= 32 {
		// generate the value
		var knuth int64 = 2654435769
		return (knuth * value) >> (32 - size)
	}

	return 0
}

// we create a simple map thing ( :) )
func (dict DictEvalNode) hashKey(key string) int64 {
	// use fibonacci hashing
	// return knuthHashing()
	key = strings.ToLower(key)
	delta := 0

	for i := 0; i < len(key); i++ {
		delta += int(key[i])
	}

	return knuthHashing(int64(delta), 32) % dict.capacity
}

func (dict DictEvalNode) resize() {
	_capacity := dict.capacity * 2
	dict.capacity = _capacity

	_dict := make([]DictItem, _capacity)

	for i := 0; i < len(dict.items); i++ {
		_dict[dict.hashKey(dict.items[i].key)] = dict.items[i]
	}

	dict.items = _dict
}

// put a value into the store
func (dict DictEvalNode) Put(key string, value interface{}) ExceptionNode {

	// check if we are almost exceeding the quota if so we need to resize the dict
	if float64(dict.capacity-int64(len(dict.items))) < float64(dict.capacity)*0.75 {
		// resize the dict
		dict.resize()
	}

	dict.items[dict.hashKey(key)] = DictItem{
		key:   key,
		value: value,
	}

	return ExceptionNode{
		Type: NO_EXCEPTION,
	}
}

// get the value from the dict
func (dict DictEvalNode) Get(key string) (interface{}, ExceptionNode) {
	// throw the error later after we implement a good abstraction :)
	_dict_item := dict.items[dict.hashKey(key)]
	return _dict_item.value, ExceptionNode{
		Type: NO_EXCEPTION,
	}

	/*
		ExceptionNode{
			Type:    "AccessException",
			Message: fmt.Sprintf("%s", "The value does not exist"),
		}
	*/
}

// remove a key from the dict
func (dict DictEvalNode) Delete(key string) error {
	// index := dict.hashKey(key)
	// _dict_item := dict.items[index]

	// if _dict_item.value == nil {
	// 	return fmt.Errorf("%s", "The value does not exist")
	// }

	return nil
}
