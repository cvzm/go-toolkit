package lru

import (
	"reflect"
	"testing"
)

// 测试缓存的value类型
type String string

func (d String) Len() int {
	return len(d)
}

// Cache.Add test. 需覆盖的功能点：
// 1. 元素是否正常添加
// 2. 元素是否正常覆盖
// 3. 超出缓存大小后，最旧的元素是否还存在。
// 4. 最旧的元素删除后，其余元素个数是否正常。即：校验是否过量删除
func TestAdd(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	bytes := len(k1 + k2 + v1 + v2)
	lru := New(int64(bytes), nil)
	lru.Add(k1, String(v1+"k1"))
	lru.Add(k1, String(v1))
	if v, _ := lru.Get(k1); string(v.(String)) != v1 {
		t.Fatalf("Add 覆盖 failed")
	}
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.entryList.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

// Cache.Get test. 需覆盖的功能点：
// 1. get获取是否成功，value值是否正确
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

// Cache.RemoveOldest test. 需覆盖的功能点：
// 1. 是否正常删除了最旧的元素
func TestRemoveOldest(t *testing.T) {
	k1, k2 := "key1", "key2"
	v1, v2 := "value1", "value2"
	bytes := len(k1 + k2 + v1 + v2)
	lru := New(int64(bytes), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	if lru.entryList.Len() != 2 {
		t.Fatalf("add cache failed")
	}
	lru.RemoveOldest()
	if _, ok := lru.Get("key1"); ok || lru.entryList.Len() != 1 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

// Cache.OnEvicted test. 需覆盖的功能点：
// 1. 缓存溢出时，Cache.OnEvicted 是否被触发
// 2. Cache.OnEvicted 触发时，对应的key，value是否是正确。
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
