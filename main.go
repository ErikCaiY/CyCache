package main

import (
	"fmt"
	"sync"
	"time"
)

//import (
//	"CyCache/lru"
//)
//
//type MyInt int64
//
//func (i MyInt) Len() int {
//	return 64
//}
//
//func main() {
//	c := lru.NewCache(int64(0), nil)
//	c.Add("key1", MyInt(123))
//}

var m sync.Mutex
var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true
	m.Unlock()
}

func main() {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
