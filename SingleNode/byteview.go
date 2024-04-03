package CyCache

// ByteView 缓存类型，使用[]byte可以支持所有类型的缓存，只读。
type ByteView struct {
	b []byte
}

// Len 返回缓存长度。 实现了lru包中Value接口的Len()方法，是Value的子类
func (bv ByteView) Len() int {
	return len(bv.b)
}

// ByteSlice 缓存切片
func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

// 转为string类型
func (bv ByteView) String() string {
	return string(bv.b)
}

// 返回缓存的复制
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
