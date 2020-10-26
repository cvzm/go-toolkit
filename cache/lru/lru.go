/*
	Least Recent Used （最近最少使用）
	wiki: https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_recently_used_(LRU)
*/

package lru

import "container/list"

// lru 实现的核心结构, 使用Map+双向链表模型。访问、添加、删除的时间复杂度均为O(1)
// 注：并发访问不安全
type Cache struct {
	maxBytes     int64                    // 最大缓存byte
	currentBytes int64                    // 当前缓存byte
	entryList    *list.List               // 存放实际元素的双向链表
	cache        map[string]*list.Element // 缓存map，存放key与value对应的list指针
	// 可选，在缓存元素被清除时执行
	OnEvicted func(key string, value Value)
}

// 存储元素结构
type entry struct {
	key   string
	value Value
}

// Value 作为缓存值的结构体，只要有 Len 函数即可，用于获取value值的大小。
type Value interface {
	Len() int
}

// 获取 LRU 缓存实例
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		entryList: list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 添加|修改缓存，同时清理溢出 cache 大小的缓存元素
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.entryList.MoveToFront(ele)
		e := ele.Value.(*entry)
		c.currentBytes += int64(e.value.Len()) - int64(value.Len())
		e.value = value
	} else {
		ele := c.entryList.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.currentBytes += int64(len(key)) + int64(value.Len())
	}
	// 处理cache溢出
	for c.maxBytes > 0 && c.maxBytes < c.currentBytes {
		c.RemoveOldest()
	}
}

// 获取缓存，同时将key移至首位
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.entryList.MoveToFront(ele)
		e := ele.Value.(*entry)
		return e.value, true
	}
	return
}

// 清除最旧的缓存数据，并执行 OnEvicted 的回调函数
// RemoveOldest 不会校验缓存是否溢出
func (c *Cache) RemoveOldest() {
	ele := c.entryList.Back()
	if ele == nil {
		return
	}
	c.entryList.Remove(ele)
	e := ele.Value.(*entry)
	delete(c.cache, e.key)
	c.currentBytes -= int64(len(e.key)) + int64(e.value.Len())
	if c.OnEvicted != nil {
		c.OnEvicted(e.key, e.value)
	}

}
