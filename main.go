package main

import (
	"CyCache/lru"
)

type MyInt int64

func (i MyInt) Len() int {
	return 64
}

func main() {
	c := lru.NewCache(int64(0), nil)
	c.Add("key1", MyInt(123))
}
