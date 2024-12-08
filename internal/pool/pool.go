package pool

import (
	"sync"

	"github.com/kuekiko/NotionGO/blocks"
	"github.com/kuekiko/NotionGO/pages"
)

// RichTextPool 富文本对象池
var RichTextPool = sync.Pool{
	New: func() interface{} {
		return &blocks.RichText{}
	},
}

// BlockPool 块对象池
var BlockPool = sync.Pool{
	New: func() interface{} {
		return &blocks.Block{}
	},
}

// PagePool 页面对象池
var PagePool = sync.Pool{
	New: func() interface{} {
		return &pages.Page{}
	},
}

// BytesPool 字节缓冲池，用于 JSON 编码
var BytesPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 1024)
		return &b
	},
}

// Get 从对象池获取对象
func Get[T any](pool *sync.Pool) *T {
	return pool.Get().(*T)
}

// Put 将对象放回对象池
func Put(pool *sync.Pool, obj interface{}) {
	pool.Put(obj)
}

// GetBytes 获取字节缓冲
func GetBytes() *[]byte {
	return BytesPool.Get().(*[]byte)
}

// PutBytes 放回字节缓冲
func PutBytes(b *[]byte) {
	*b = (*b)[:0]
	BytesPool.Put(b)
}
