# 单机并发缓存
之前实现的lru算法在并发条件下并不安全，需要针对并发机制做优化

## ByteView
构造一个数据结构ByteView实现Value

ByteView需要实现Len()方法
- Len() int
- ByteSlice() []byte // 获得一份副本
- String() string // 获得缓存的string类型副本
- cloneBytes(b []byte) []byte // 获得一份副本(内部方法)

## cache封装lru
加入sync.Mutex实现线程安全

- add(key string, value ByteView)
- get(key string) (value ByteView, ok bool)

## 主体结构Group
Group用于与用户的交互

Group主要思路如下
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