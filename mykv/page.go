package mykv

import (
	"unsafe"
	"hash/fnv"
	"fmt"
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
	// 用于标记页头结尾处，或者页面内存储数据的起始处。
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

//一个branchPage或leafPage由页头和若干branchPageElements或leafPageElements组成

type branchPageElement struct {
	//element对应的K/V对存储位置相对于当前element的偏移
	pos   uint32
	//element对应的Key的长度，以字节为单位
	ksize uint32
	//element指向的子节点所在page的页号
	pgid  pgid
}

type leafPageElement struct {
	//标明当前element是否代表一个Bucket，如果是Bucket则其值为1，如果不是则其值为0
	flags uint32
	pos   uint32
	ksize uint32
	vsize uint32
}

func (p *page)meta()*meta{
	return (*meta)(unsafe.Pointer(&p.ptr))
}

func (m *meta) sum64() uint64 {
	var h = fnv.New64a()
	h.Write((*[unsafe.Offsetof(meta{}.chckksum)]byte)(unsafe.Pointer(m))[:])
	return h.Sum64()
}
// 验证mata page的有效性
func (m *meta) validate() error {
	if m.magic != magic {
		return ErrInvalid
	}else if m.version != version {
		return ErrVersionMismatch
	}else if m.chckksum != 0 && m.chckksum != m.sum64(){
		return ErrChecksum
	}
	return nil
}
// meta 页的存储
func (m *meta) write (p *page) {
	if m.root.root >= m.pgid{
		fmt.Println("here1")
	}else if m.freelist >= m.pgid {
		fmt.Println("here2")
	}
//	page id 可以是0或者1，我们可以通过交易ID来确定
	p.id = pgid(m.txid%2)
	p.flags |= metaPageFlag

	m.chckksum = m.sum64()
	m.copy(p.meta())
}
func (m *meta) copy (dest *meta){
	*dest = *m
}