/*
	All right reserved：https://github.com/hunterhug/gomap at 2020
	Attribution-NonCommercial-NoDerivatives 4.0 International
	You can use it for education only but can't make profits for any companies and individuals!
*/
package gomap // import "github.com/hunterhug/gomap"

import "strings"

// Map method
// design to be concurrent safe
// should support int key?
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
	Check() bool                                                  // just help
}

// Iterator concurrent not safe
// you should deal by yourself
type MapIterator interface {
	HasNext() bool
	Next() (key string, value interface{})
}

// default map is rbt implement
func NewMap() Map {
	return &rbTree{}
}

// default map is rbt implement
func New() Map {
	return &rbTree{}
}

// compare two key
func compare(key1, key2 string) int64 {
	return int64(strings.Compare(key1, key2))
}
