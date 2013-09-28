package storage

import (
	"container/list"
	"sync"
)

// 线程安全的List，提供满足Redis List的函数
type SafeList struct {
	innerList *list.List
	mutex     *sync.Mutex
}

func NewSafeList() (sl *SafeList) {
	sl = &SafeList{}
	sl.innerList = list.New()
	sl.mutex = &sync.Mutex{}
	return
}

func (sl *SafeList) LPop() (value interface{}) {
	sl.mutex.Lock()
	elem := sl.innerList.Front()
	if elem != nil {
		value = elem.Value
		sl.innerList.Remove(elem)
	}
	sl.mutex.Unlock()
	return
}

func (sl *SafeList) RPop() (value interface{}) {
	sl.mutex.Lock()
	elem := sl.innerList.Back()
	if elem != nil {
		value = elem.Value
		sl.innerList.Remove(elem)
	}
	sl.mutex.Unlock()
	return
}

func (sl *SafeList) LPush(values ...string) (length int) {
	sl.mutex.Lock()
	for _, value := range values {
		sl.innerList.PushFront(value)
	}
	length = sl.innerList.Len()
	sl.mutex.Unlock()
	return
}

func (sl *SafeList) RPush(values ...string) (length int) {
	sl.mutex.Lock()
	for _, value := range values {
		sl.innerList.PushBack(value)
	}
	length = sl.innerList.Len()
	sl.mutex.Unlock()
	return
}

func (sl *SafeList) Len() (length int) {
	sl.mutex.Lock()
	length = sl.innerList.Len()
	sl.mutex.Unlock()
	return
}

// 枚举实现，超大列表下性能不佳，并且lock住其它操作
func (sl *SafeList) Index(index int) (value interface{}) {
	sl.mutex.Lock()
	i := 0
	for e := sl.innerList.Front(); e != nil; e = e.Next() {
		if i == index {
			value = e.Value
			break
		}
		i++
	}
	sl.mutex.Unlock()
	return
}

// 枚举实现，超大列表下性能不佳，并且lock住其它操作
func (sl *SafeList) Range(start int, end int) (values []interface{}) {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	length := sl.innerList.Len()
	if start >= length || end < start {
		values = make([]interface{}, 0)
		return
	}
	// 确认返回数组大小
	resultsize := 0
	if length > end {
		resultsize = end - start + 1
	} else {
		resultsize = length - start
	}
	values = make([]interface{}, 0, resultsize)
	// 填充数据
	i := 0
	for e := sl.innerList.Front(); e != nil; e = e.Next() {
		if i >= start && i <= end {
			values = append(values, e.Value)
		}
		i++
	}
	return
}
