# Map Golang implement by Red-Black Tree, AVL Tree

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gomap.svg?style=social&label=Forks)](https://github.com/hunterhug/gomap/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gomap.svg?style=social&label=Stars)](https://github.com/hunterhug/gomap/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap/issues)

[中文说明](/README_ZH.md)

Map implement by tree data struct such Red-Black Tree, AVL Tree.

Our tree map is design to be concurrent safe, and don't waste any space different from golang standard map which not shrink even map key pairs num is 0.

## Usage

Simple get it by:

```
go get -v github.com/hunterhug/gomap
```

There are:

1. Standard Red-Black Tree Map(2-3-4-Tree): `gomap.New()`，`gomap.NewMap()`,`gomap.NewRBMap()`.
2. AVL Tree Map: `gomap.NewAVLMap()`.

Core api:

```go
// Map method
// design to be concurrent safe
type Map interface {
	Put(key string, value interface{})                            // put key pairs
	Delete(key string)                                            // delete a key
	Get(key string) (value interface{}, exist bool)               // get value from key
	GetInt(key string) (value int, exist bool, err error)         // get value auto change to Int
	GetInt64(key string) (value int64, exist bool, err error)     // get value auto change to Int64
	GetString(key string) (value string, exist bool, err error)   // get value auto change to string
	GetFloat64(key string) (value float64, exist bool, err error) // get value auto change to string
	GetBytes(key string) (value []byte, exist bool, err error)    // get value auto change to []byte
	Contains(key string) (exist bool)                             // map contains key?
	Len() int64                                                   // map key pairs num
	KeyList() []string                                            // map key out to list from top to bottom which is layer order
	KeySortedList() []string                                      // map key out to list sorted
	Iterator() MapIterator                                        // map iterator, iterator from top to bottom which is layer order
	MaxKey() (key string, value interface{}, exist bool)          // find max key pairs
	MinKey() (key string, value interface{}, exist bool)          // find min key pairs
	SetComparator(comparator) Map                                 // set compare func to control key compare
}

// Iterator concurrent not safe
// you should deal by yourself
type MapIterator interface {
	HasNext() bool
	Next() (key string, value interface{})
}
```

We has already implement them by non recursion way and optimized a lot, so use which type of tree map is no different.

## Example

Some example below:

```go
package main

import (
	"fmt"
	"github.com/hunterhug/gomap"
	"math/rand"
	"strconv"
	"time"
)

// loop times
// 循环的次数，用于随机添加和删除键值对
// 可以修改成1000
var num = 20

func init() {
	// random seed
	// 随机数种子初始化
	rand.Seed(time.Now().Unix())
}

// diy comparator
// 内部是按键字符串来做查找树，我们把他变成整形
func comparatorInt(key1, key2 string) int64 {
	k1, _ := strconv.Atoi(key1)
	k2, _ := strconv.Atoi(key2)
	if k1 == k2 {
		return 0
	} else if k1 > k2 {
		return 1
	} else {
		return -1
	}
}

func main() {
	checkMap := make(map[string]struct{})

	// 1. new a map default is rb tree
	// 1. 新建一个 Map，默认为标准红黑树实现
	m := gomap.New()
	//m = gomap.NewAVLMap()    // avl tree better version
	//m = gomap.NewAVLRecursionMap() // avl tree bad version

	m.SetComparator(comparatorInt) // set inner comparator

	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("add key:", key)
		checkMap[key] = struct{}{}

		// 2. put key pairs
		// 2. 放键值对
		m.Put(key, key)
	}

	fmt.Println("map len is ", m.Len())

	// 3. can iterator
	// 3. 迭代器使用
	iterator := m.Iterator()
	for iterator.HasNext() {
		k, v := iterator.Next()
		fmt.Printf("Iterator key:%s,value %v\n", k, v)
	}

	// 4. get key
	// 4. 获取键中的值
	key := "9"
	value, exist := m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 5. get int will err
	// 5. 获取键中的值，并且指定类型，因为类型是字符串，所以转成整形会报错
	_, _, err := m.GetInt(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 6. check is a rb tree
	// 6. 内部辅助，检查是不是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	// 7. delete '9' then find '9'
	// 7. 删除键 '9' 然后再找 '9'
	m.Delete(key)
	value, exist = m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 8. key list
	// 8. 获取键列表
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 9. delete many
	// 9. 删除很多键值对
	for key := range checkMap {
		v, ok := m.Get(key)
		fmt.Printf("find %s:%v=%v, delete key:%s\n", key, v, ok, key)
		m.Delete(key)
		if !m.Check() {
			fmt.Println("is not a rb tree,len:", m.Len())
		}
	}

	// 10. key list
	// 10. 获取键列表
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 11. check is a rb tree
	// 11. 再次检查是否是一颗正常的红黑树
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
```

## BenchTest

We test `Golang map` and `Red-Black Tree`, `AVL Tree`:

```go
go test -run="bench_test.go" -test.bench=".*" -test.benchmem=1 -count=1                                                                                                            master 
goos: darwin
goarch: amd64
pkg: github.com/hunterhug/gomap
BenchmarkGolangMapPut-12                 1791264               621 ns/op             112 B/op          6 allocs/op
BenchmarkRBTMapPut-12                    1000000              1408 ns/op             104 B/op          6 allocs/op
BenchmarkAVLMapPut-12                    1000000              1440 ns/op             104 B/op          6 allocs/op
BenchmarkAVLRecursionMapPut-12           1000000              2024 ns/op             104 B/op          6 allocs/op
BenchmarkGolangMapDelete-12              3577232               303 ns/op              15 B/op          1 allocs/op
BenchmarkRBTMapDelete-12                  996924              1248 ns/op              15 B/op          1 allocs/op
BenchmarkAVLMapDelete-12                 1000000              1227 ns/op              15 B/op          1 allocs/op
BenchmarkAVLRecursionMapDelete-12         667242              1866 ns/op              15 B/op          1 allocs/op
BenchmarkGolangMapGet-12                15797131                72.2 ns/op             2 B/op          1 allocs/op
BenchmarkRBTMapGet-12                    5798295               195 ns/op               2 B/op          1 allocs/op
BenchmarkAVLMapGet-12                    5831353               197 ns/op               2 B/op          1 allocs/op
BenchmarkAVLRecursionMapGet-12           5275490               228 ns/op               2 B/op          1 allocs/op
BenchmarkGolangMapRandom-12              1256779               940 ns/op             146 B/op          8 allocs/op
BenchmarkRBTMapRandom-12                  965804              2652 ns/op             126 B/op          8 allocs/op
BenchmarkAVLMapRandom-12                  902004              2805 ns/op             126 B/op          8 allocs/op
BenchmarkAVLRecursionMapRandom-12         701880              3309 ns/op             129 B/op          8 allocs/op
PASS
ok      github.com/hunterhug/gomap      66.006s
```

If you want to save memory space, you can choose our tree map.

# License

```
Copyright [2019-2021] [github.com/hunterhug]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```