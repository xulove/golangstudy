package mykv

import (
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"
)

// IgnoreNoSync specifies（指定） whether the NoSync field of a DB is ignored when
// syncing changes to a file.  This is required as some operating systems（操作系统）,
// such as OpenBSD, do not have a unified buffer cache (UBC) and writes
// must be synchronized using the msync(2) syscall.

// Options represents the options that can be set when opening a database.
type Options struct {
	// Timeout is the amount of time to wait to obtain（获得） a file lock.
	// When set to zero it will wait indefinitely（无限期）. This option is only
	// available on Darwin and Linux.
	Timeout time.Duration

	// Sets the DB.NoGrowSync flag before memory mapping the file.
	NoGrowSync bool

	// Open database in read-only mode. Uses flock(..., LOCK_SH |LOCK_NB) to
	// grab a shared lock (UNIX).
	ReadOnly bool

	// Sets the DB.MmapFlags flag before memory mapping the file.
	MmapFlags int

	// InitialMmapSize is the initial mmap size of the database
	// in bytes. Read transactions won't block write transaction
	// if the InitialMmapSize is large enough to hold database mmap
	// size. (See DB.Begin for more information)
	//
	// If <=0, the initial map size is 0.
	// If initialMmapSize is smaller than the previous database size,
	// it takes no effect.
	InitialMmapSize int
}

// DefaultOptions represent the options used if nil options are passed into Open().
// No timeout is used which will cause Bolt to wait indefinitely for a lock.
var DefaultOptions = &Options{
	Timeout:    0,
	NoGrowSync: false,
}

// Default values if not set in a DB instance.
const (
	DefaultMaxBatchSize  int = 1000
	DefaultMaxBatchDelay     = 10 * time.Millisecond
	DefaultAllocSize         = 16 * 1024 * 1024
)

type DB struct {
	//启用后，数据库每次提交都会进行一次check（）
	//当数据库状态不一致时，就会报错
	//strickMode对性能影响很大，建议在调试模式下才使用
	StrictMode bool
	// Setting the NoSync flag will cause the database to skip fsync() calls after each commit.
	// This can be useful when bulk loading data(批量加载数据) into a database
	// and you can restart the bulk load in the event of a system failure or database corruption（损坏）.
	// Do not set this flag for normal use.
	// If the package global IgnoreNoSync constant is true, this value is
	// ignored.  See the comment on that constant for more details.
	NoSybc bool
	// When true, skips the truncate call when growing the database.
	// Setting this to true is only safe on non-ext3/ext4 systems.
	// Skipping truncation avoids preallocation（预分配） of hard drive space and
	// bypasses a truncate() and fsync() syscall on remapping.

	NoGrowSync bool
	// If you want to read the entire database fast, you can set MmapFlag to
	// syscall.MAP_POPULATE on Linux 2.6.23+ for sequential read-ahead(顺序预读).
	MmapFlags     int
	MaxBatchSize  int
	MaxBatchDelay time.Duration
	AllocSize     int

	path     string
	file     *os.File
	lockfile *os.File // windows only
	//这三个带data的字段应该和内存映射有关系
	dataref  []byte // mmap'ed readonly, write throws SEGV
	data     *[maxMapSize]byte
	datasz   int
	filesz   int // current on disk file size
	meta0    *meta
	meta1    *meta
	pageSize int
	opened   bool
	//rwtx     *Tx
	//txs      []*Tx
	freelist *freelist
	//stats    Stats

	pagePool sync.Pool

	batchMu sync.Mutex
	//batch   *batch

	rwlock   sync.Mutex   // Allows only one writer at a time.
	metalock sync.Mutex   // Protects meta page access.
	mmaplock sync.RWMutex // Protects mmap access during remapping.
	statlock sync.RWMutex // Protects stats access.

	ops struct {
		writeAt func(b []byte, off int64) (n int, err error)
	}

	// Read only mode.
	// When true, Update() and Begin(true) return ErrDatabaseReadOnly immediately.
	readOnly bool
}

func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	// 创建数据库，并将状态设置为opend
	var db = &DB{opened: true}
	if options == nil {
		options = DefaultOptions
	}
	db.NoGrowSync = options.NoGrowSync
	db.MmapFlags = options.MmapFlags
	// Set default values for later DB operations.
	db.MaxBatchSize = DefaultMaxBatchSize
	db.MaxBatchDelay = DefaultMaxBatchDelay
	db.AllocSize = DefaultAllocSize
	flag := os.O_RDWR
	if options.ReadOnly {
		flag = os.O_RDONLY
		db.readOnly = true
	}
	// Open data file and separate sync handler for metadata writes.
	db.path = path
	var err error
	if db.file, err = os.OpenFile(path, flag|os.O_CREATE, mode); err != nil {
		//_ = db.close()
		return nil, err
	}
	if err := flock(db, mode, !db.readOnly, options.Timeout); err != nil {
		//_ = db.Close()
		return nil, err
	}
	db.ops.writeAt = db.file.WriteAt

	if info, err := db.file.Stat(); err != nil {
		return nil, err
	} else if info.Size() == 0 {
		if err := db.init(); err != nil {
			return nil, err
		}
	} else {
		// read the first meta page tp determine the page size
		var buf [0x1000]byte
		if _,err := db.file.WriteAt(buf[:],0);err == nil {
			m := db.pageInBuffer(buf[:],0).meta()
			if err := m.validate();err != nil {
				db.pageSize = os.Getpagesize()
			}else{
				db.pageSize = int(m.pageSize)
			}
		}
	}
	//options.InitialMmapSize此时应该是0
	if err := db.mmap(options.InitialMmapSize);err != nil {
		_ = db.close()
		return nil,err
	}

}

// create a new database file and initializes its meta page
func (db *DB) init() error {
	db.pageSize = os.Getpagesize()

	buf := make([]byte, db.pageSize*4)
	//将第0页和第1页初始化meta页，
	// 并指定root bucket的page id为3, freelist记录的page id为2，
	// 当前数据库总页数为4，同时txid分别为0和1
	for i := 0; i < 2; i++ {
		p := (*page)(unsafe.Pointer(&buf[pgid(i)*pgid(db.pageSize)]))
		p.id = pgid(i)
		p.flags = metaPageFlag

		// initialize the meta page
		m := p.meta()

		m.magic = magic
		m.version = version
		m.pageSize = uint32(db.pageSize)
		m.freelist = 2
		m.root = bucket{root: 3}
		m.pgid = 4

		m.txid = txid(i)
		m.chckksum = m.sum64()
	}
	// 将第2页初始化为freelist页，即freelist的记录将会存在第2页；
	p := db.pageInBuffer(buf[:], pgid(2))
	p.id = pgid(2)
	p.flags = freelistPageFlag
	p.count = 0
	// 将第3页初始化为一个空页，它可以用来写入K/V记录，请注意它必须是B+ Tree中的叶子节点
	p = db.pageInBuffer(buf[:], pgid(3))
	p.id = pgid(3)
	p.flags = leafPageFlag
	p.count = 0
	//调用写文件函数将buffer中的数据写入文件
	// 在open函数中已经设定 db.ops.writeAt = db.file.writeAt.
	// 所以现在写数据就是往db.file中写数据
	if _, err := db.ops.writeAt(buf, 0); err != nil {
		return err
	}
	//通过fdatasync()调用将内核中磁盘页缓冲立即写入磁盘
	if err := fdatasync(db); err != nil {
		return err
	}
	return nil
}
func (db *DB) mmap(minsz int) error {
	db.mmaplock.Lock()
	defer db.mmaplock.Unlock()
	info, err := db.file.Stat()
	if err != nil {
		return fmt.Errorf("mmap stat error:%s", err)
	} else if int(info.Size()) < db.pageSize*2 {
		return fmt.Errorf("file size is too small")
	}
	var size = int(info.Size())
	if size < minsz {
		size = minsz
	}
	// 传入的size是16kb，穿出的返会的size已经是2^15,是32kb了
	size, err = db.mmapSize(size)
	if err != nil {
		return err
	}
	// 如果内存中已经映射了，清除映射
	if err := db.munmap(); err != nil {
		return err
	}
	// 真正的开始映射
	if err := mmap(db, size); err != nil {
		return nil
	}
	db.meta0 = db.page(0).meta()
	db.meta1 = db.page(1).meta()
	// 验证meta page的有效性。
	// 只有当两者都失效的情况下，我们才会错误
	// 一个失效的话，我们可以从另一个进行恢复
	err0 := db.meta0.validate()
	err1 := db.meta1.validate()
	if err0 != nil && err1 != nil{
		return err0
	}
	return nil

}

// 1 kb = 1024 bytes =2^10 bytes
// 1 Mb = 1024 kb = 2^20 bytes
// 1 Gb = 1024 Mb = 2^30 bytes
// 1<<3，相当于1*2^3
func (db *DB) mmapSize(size int) (int, error) {
	for i := uint(15); i <= 30; i++ {
		if size <= 1<<i {
			return 1 << i, nil
		}
	}
	if size > maxMapSize {
		return 0, fmt.Errorf("mmap too large")
	}
	sz := int64(size)
	//如果大于1GB则每次增长1GB
	if reminder := sz % int64(maxMmapStep); reminder > 0 {
		sz += int64(maxMmapStep) - reminder
	}
	//确保mmap的大小是pageSize的整数倍
	pageSize := int64(db.pageSize)
	if (sz % pageSize) != 0 {
		sz = ((sz / pageSize) + 1) * pageSize
	}
	if sz > maxMapSize {
		sz = maxMapSize
	}
	return int(sz), nil
}
func (db *DB) munmap() error {
	if err := munmap(db); err != nil {
		return fmt.Errorf("unmap error:" + err.Error())
	}
	return nil
}
func (db *DB) pageInBuffer(b []byte, id pgid) *page {
	return (*page)(unsafe.Pointer(&b[id*pgid(db.pageSize)]))
}

func (db *DB) page(id pgid) *page {
	pos := id * pgid(db.pageSize)
	return (*page)(unsafe.Pointer(&db.data[pos]))
}

func (db *DB) close() error {
	if !db.opened {
		return nil
	}
	db.opened = false
	db.freelist = nil

	db.ops.writeAt = nil

	if err := db.munmap();err != nil {
		return err
	}
	// close file handles
	if db.file != nil {
		// no need to unlock read-only file
		if !db.readOnly {
			// 解锁文件，暂时先不用实现了
		}

		// 关闭文件描述符,调用golang的文件关闭函数
		if err := db.file.Close();err != nil {
			return fmt.Errorf("db file close:%s",err)
		}

		db.file = nil
	}
	db.path = ""
	return nil

}