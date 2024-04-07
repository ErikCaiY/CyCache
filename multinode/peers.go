package CyCache

type PeerPicker interface {
	// PickPeer 用于选择相应的节点
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	// Get 从对应的group查找缓存值
	Get(group string, key string) ([]byte, error)
}
