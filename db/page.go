package db

import (
	"unsafe"
	"hash/fnv"
)
//四种页面类型
const branchPageFlag  = 0x01  //对应B+ tree中的内节点
const leafPageFlag  = 0x02   //对应B+ tree中的叶子节点
const metaPageFlag = 0x04
const freelistPageFlag  = 0x10

type pgid uint64
// 一个page页面=page header（没有ptr） + elements
type page struct {
	// 页面的id，如0,1,2.从数据库文件内存映射中读取一页的索引值
	id pgid
	// 页面类型
	flags uint16
	// 页面内存储的元素个数。只在branchPage中leafPage有用。
	// 对应的元素分别是branchPageElement和leafPageElement
	count uint16
	// 当前页是否有后续页。
	// 如果有，表示后续页的数量。如果没有，则为0
	overflow uint64
	// 用于标记也头结尾处，或者业内存储数据的起始处。
	ptr uintptr
}
// meta page的结构
type meta struct {
	//boltdb的magic number
	magic uint32
	// boltdb的version
	version uint32
	// boltdb的页的大小
	pageSize uint32
	// flags 保留字段，暂时没有用到
	flags uint32
	// root :boltdb根bucket的头信息
	root bucket
	// boltdb中存freelist的页号。freelist用来存空闲页面的页号
	freelist pgid
	// pgid 简单理解成boltdb文件中的总页数
	pgid pgid
	// 上一次写数据库的transcation id
	txid txid
	// checksum 上面所有字段64位哈希校验
	chckksum uint64
}
func (p *page)meta()*meta{
	return (*meta)(unsafe.Pointer(&p.ptr))
}

func (m *meta) sum64() uint64 {
	var h = fnv.New64a()
	h.Write((*[unsafe.Offsetof(meta{}.chckksum)]byte)(unsafe.Pointer(m))[:])
	return h.Sum64()
}