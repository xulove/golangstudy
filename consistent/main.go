package main
import "hash/crc32"
import "strconv"
import (
	"errors"
	"sort"
	"sync"
	"fmt"
)

type HashRing []uint32
const default_replicas = 160

func (h HashRing) Len() int {
	return len(h)
}

func (h HashRing) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h HashRing) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

type Consistent struct {
	Nodes map[uint32]Node // map的key就是把Node各个字段hash之后的结果，uint32类型
	numReps int   //相当于把物理节点虚拟化，本来只有一个物理节点，通过这个参数就可以虚拟出160个节点
	Resources map[int]bool  //用node的id作为map的key，这就是用来判断某个node是否存在在此分布式系统中
	ring HashRing
	sync.RWMutex
}

type Node struct {
	Id int
	Ip string
	Port int
	HostName string
	Weight int //这是我们引入了一个权重的概念
}

func NewConsistent ()*Consistent{
	return &Consistent{
		Nodes:make(map[uint32]Node),
		numReps:default_replicas,
		Resources:make(map[int]bool),
		ring:HashRing{},
	}
}
func NewNode(id int,ip string,port int,hostName string,weight int)*Node{
	return &Node{ id,ip,port,hostName,weight}
}
//添加一个节点
func (c *Consistent) Add (node *Node) error{
	c.Lock()
	defer c.Unlock()
	
	if _,ok := c.Resources[node.Id];ok {
		return errors.New("node has existed")
	}
	// 如果你的weight高，那么就多给你几个i，然你多生成几个hash。
	// 如果没有i，每次生成的大小不变，导致在圆环中的位置也不变，那就没有意思了。此处的i很像比特币挖矿处的nonce值
	count := c.numReps * node.Weight
	for i:=0;i<count;i++ {
		str := joinStr(i,node)
		c.Nodes[hashStr(str)] = *node
	}
	c.Resources[node.Id] = true
	c.sortHashRing()
	return nil
}
// 对*Consistent系统重新进行排序，排序的依据就是Nodes字段map中的key，ring中存储的map字段的key。也是每个节点的hash后的大小，标识
func (c *Consistent) sortHashRing(){
	c.ring = []uint32{}
	for k,_ := range c.Nodes {
		c.ring = append(c.ring,k)
	}
	sort.Sort(c.ring)
}

//这里就是你给我一个key，我帮你找到你该存储到哪个节点Node上
func (c *Consistent) Get(key string) Node {
	c.Lock()
	defer c.Unlock()
	
	hash := hashStr(key)
	// 根据hash找该存储到环上的第i个节点上
	i := c.search(hash)
	return c.Nodes[c.ring[i]]	
}
//根据给定的hash，找此hash对应的数据应该存储到一致性环ring中的第几个节点上
// 什么hash该存储在什么样的节点上，可以根据自己的规则来
func (c *Consistent) search(hash uint32) int {
	i := sort.Search(len(c.ring),func(i int)bool{ return c.ring[i]>=hash})
	if i == len(c.ring){
		return 0
	}
	return i
}
func (c *Consistent) Remove(node *Node) error {
	c.Lock()
	defer c.Unlock()
	
	if _,ok := c.Resources[node.Id];!ok {
		return errors.New("node not exist")
	}
	count := c.numReps * node.Weight
	for i:=0;i<count;i++ {
		str := joinStr(i,node)
		delete(c.Nodes,hashStr(str))
	}
	c.sortHashRing()
	return nil
}


//hash算法
func hashStr (key string) uint32{
	return crc32.ChecksumIEEE([]byte(key))
}
// 合并node中的一些字段
func joinStr (i int,node *Node)string{
	return node.Ip + "*" + strconv.Itoa(node.Weight) + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(node.Id)
}

func main(){
	ins := NewConsistent()
	for i:=0;i<10;i++ {
		si := fmt.Sprintf("%d",i)
		ins.Add(NewNode(i,"127.18.1."+si,8090,"host_"+si,1))
	}
	for k,v := range ins.Nodes{
		fmt.Println("Hash:",k, "IP:",v.Ip)
	}
	ipMap := make(map[string]int,0)
	for i :=0;i<10000;i++{
		si := fmt.Sprintf("key%d",i)
		n := ins.Get(si)
		if _,ok := ipMap[n.Ip];ok {
			ipMap[n.Ip] += 1
		}else{
			ipMap[n.Ip] = 1
		}
	}
	for k, v := range ipMap {
		fmt.Println("Node IP:", k, " count:", v)
	}

	
}
