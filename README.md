# CyCache
一个基于Go语言以及LRU算法的缓存管理系统

## LRU算法 
基于map与双向链表实现

链表头部表示最近使用过的元素

链表尾部表示最近最少使用的元素(优先删除)

考虑到实际使用时可能会涉及到不同的数据类型，因此构建一个Value接口表示数据类型，
所有使用的数据类型都要实现Value接口进行使用

#### LRU提供的方法
- Get(key string) (value Value, ok bool)
- RemoveOldest() error
- ADD(key string, value Value)

## ByteView
构造一个数据结构ByteView描述缓存存储的内容

ByteView需要实现Len()方法
- Len() int
- ByteSlice() []byte // 获得一份副本
- String() string // 获得缓存的string类型副本
- cloneBytes(b []byte) []byte // 获得一份副本(内部方法)

## Cache类对实现的LRU算法进行封装，加入线程锁保证线程安全

- add(key string, value ByteView)
- get(key string) (value ByteView, ok bool)

## 主题结构Group用于用户的交互
```

接收一个key ---> 是否命中缓存 ---Yes---> 返回缓存值
            |
            No
            |---> 是否从本地获取 ---Yes---> (调用回调函数),获取值并添加到缓存 ---> 返回缓存
                        |
                        No
                        |---> 与远程节点交互，获取缓存值 ---> 返回缓存 
```

Group 结构体中包含
- name 一个string，表示当前Group的编号信息
- getter 一个函数式接口(既可以当接口用，也可以当函数用)
- mainCache 一个cache，表示当前Group对应存储的缓存
* Get(key string) (ByteView, error) 尝试从缓存中获取想要的值，若缓存未命中，则本地获取
    * load(key string) (value ByteView, err error)
    * getLocally(key string) (ByteView, error)
    * populateCache(key string, value ByteView)
