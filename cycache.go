package CyCache

import (
	pb "CyCache/cycachepb"
	"errors"
	"log"
	"sync"
)

// Getter 提供给外部的接口，根据键获取缓存中的值
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 用函数方法实现Getter
type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

// ==========================================================================

// 提供给外界调用的类型，一个Group可以被认为是一个缓存的命名空间
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取源数据
	mainCache cache  // 并发缓存
	peers     PeerPicker
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group) // 记录存在的group的全局变量
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if len(key) == 0 {
		return ByteView{}, errors.New("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[key] %s hits", key)
		return v, nil
	}

	// 缓存不存在，调用回调函数获取源数据添加到缓存
	return g.load(key)
}

//func (g *Group) load(key string) (value ByteView, err error) {
//	return g.getLocally(key)
//}

// 从本地获取
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 将本次未命中的值加入Group缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterPeers 实现了 PeerPicker 接口的 HTTPPool 注入到 Group 中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// load 使用 PickPeer() 方法选择节点，若非本机节点，则调用 getFromPeer() 从远程获取。若是本机节点或失败，则回退到 getLocally()
func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}

	return g.getLocally(key)
}

// 使用实现了 PeerGetter 接口的 httpGetter 从访问远程节点，获取缓存值
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
