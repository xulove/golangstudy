package mykv

import "unsafe"

type node struct {
	isLeaf bool
	inodes inodes
}
type inode struct {
	flag  uint32
	pgid  pgid
	key   []byte
	value []byte
}
type inodes []inode

func (n *node) write(p *page) {
	// byte数组指针b，它指向p.ptr指向的位置向后偏移n.pageElementSize()*len(n.inodes)的位置。
	// 上文中我们介绍过p.ptr实际上指向页头结尾处或者页正文开始处，所以b实际上是指向了页正文中elements的结尾处
	// 数组的指针就是数据中第一个元素的指针
	//[n.pageElementSize()*len(n.inodes):],没有后面这一段，前面就是个数组指针，加上这个直接就变成单纯的切片了
	b := (*[maxAllocSize]byte)(unsafe.Pointer(&p.ptr))[n.pageElementSize()*len(n.inodes):]

}

func (n *node) pageElementSize() int {
	if n.isLeaf {
		return leafPageElementSize
	}
	return branchPageElementSize
}
