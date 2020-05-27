# Map implement by Black-Red Tree, AVL Tree

[![GitHub forks](https://img.shields.io/github/forks/hunterhug/gomap.svg?style=social&label=Forks)](https://github.com/hunterhug/gomap/network)
[![GitHub stars](https://img.shields.io/github/stars/hunterhug/gomap.svg?style=social&label=Stars)](https://github.com/hunterhug/gomap/stargazers)
[![GitHub last commit](https://img.shields.io/github/last-commit/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap)
[![Go Report Card](https://goreportcard.com/badge/github.com/hunterhug/gomap)](https://goreportcard.com/report/github.com/hunterhug/gomap)
[![GitHub issues](https://img.shields.io/github/issues/hunterhug/gomap.svg)](https://github.com/hunterhug/gomap/issues)


Map implement by tree data struct such Black-Red Tree, AVL Tree.

Our tree map is design to be concurrent safe, and don't waste any space different from golang standard map which not shrink even map key pairs num is 0.

## Usage

simple get it by:

```
go get -v github.com/hunterhug/gomap
```

core api:

```go
// Map method
// design to be concurrent safe
type Map interface {
	Put(key string, value interface{})            // put key pairs
	Delete(key string)                            // delete a key
	Get(key string) (interface{}, bool)           // get value from key
	GetInt(key string) (int, bool, error)         // get value auto change to Int
	GetInt64(key string) (int64, bool, error)     // get value auto change to Int64
	GetString(key string) (string, bool, error)   // get value auto change to string
	GetFloat64(key string) (float64, bool, error) // get value auto change to string
	GetBytes(key string) ([]byte, bool, error)    // get value auto change to []byte
	Contain(key string) bool                      // map contains key?
	Len() int64                                   // map key pairs num
	KeyList() []string                            // map key out to list from top to bottom which is layer order
	KeySortedList() []string                      // map key out to list sorted
	Iterator() MapIterator                        // map iterator, iterator from top to bottom which is layer order
}

// Iterator concurrent not safe
// you should deal by yourself
type MapIterator interface {
	HasNext() bool
	Next() (key string, value interface{})
}
```

## Example

```go
package main

import (
	"fmt"
	"github.com/hunterhug/gomap"
	"math/rand"
	"time"
)

// loop times
var num = 10

func init() {
	// random seed
	rand.Seed(time.Now().Unix())
}

func main() {
	// 1. new a map
	m := gomap.New()
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("add key:", key)
		// 2. put key pairs
		m.Put(key, key)
	}

	fmt.Println("map len is ", m.Len())

	// 3. can iterator
	iterator := m.Iterator()
	for iterator.HasNext() {
		k, v := iterator.Next()
		fmt.Printf("Iterator key:%s,value %v\n", k, v)
	}

	// 4. get key
	key := "9"
	value, exist := m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 5. get int will err
	_, _, err := m.GetInt(key)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 6. check is a rb tree
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}

	// 7. delete '9' then find '9'
	m.Delete(key)
	value, exist = m.Get(key)
	if exist {
		fmt.Printf("%s exist, value is:%s\n", key, value)
	} else {
		fmt.Printf("%s not exist\n", key)
	}

	// 8. delete many
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("%d", rand.Int63n(int64(num)))
		fmt.Println("delete key:", key)
		m.Delete(key)
	}

	// 9. key list
	fmt.Printf("keyList:%#v,len:%d\n", m.KeyList(), m.Len())
	fmt.Printf("keySortList:%#v,len:%d\n", m.KeySortedList(), m.Len())

	// 10. check is a rb tree
	if m.Check() {
		fmt.Println("is a rb tree,len:", m.Len())
	}
}
```