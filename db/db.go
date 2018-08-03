package db

import (
	"os"
	"runtime"
	"time"
	"sync"
	"flag"
)
// IgnoreNoSync specifies（指定） whether the NoSync field of a DB is ignored when
// syncing changes to a file.  This is required as some operating systems（操作系统）,
// such as OpenBSD, do not have a unified buffer cache (UBC) and writes
// must be synchronized using the msync(2) syscall.

// runtime包
//尽管 Go 编译器产生的是本地可执行代码，
// 这些代码仍旧运行在 Go 的 runtime（这部分的代码可以在 runtime 包中找到）当中。
// 这个 runtime 类似 Java 和 .NET 语言所用到的虚拟机，
// 它负责管理包括内存分配、垃圾回收（第 10.8 节）、栈处理、goroutine、channel、切片（slice）、map 和反射（reflection）等等
const IgnoreNoSync = runtime.GOOS == "openbsd"
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
	MmapFlags int
	MaxBatchSize int
	MaxBatchDelay time.Duration
	AllocSize int

	path     string
	file     *os.File
	lockfile *os.File // windows only
	dataref  []byte   // mmap'ed readonly, write throws SEGV
	data     *[maxMapSize]byte
	datasz   int
	filesz   int // current on disk file size
	meta0    *meta
	meta1    *meta
	pageSize int
	opened   bool
	rwtx     *Tx
	txs      []*Tx
	freelist *freelist
	stats    Stats

	pagePool sync.Pool

	batchMu sync.Mutex
	batch   *batch

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

func open(path string,mode os.FileMode,options *Options)(*DB,error){
	// 创建数据库，并将状态设置为opend
	var db = &DB{opened:true}
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
		db.readOnly= true
	}
	// Open data file and separate sync handler for metadata writes.
	db.path = path
	var err error
	if db.file,err = os.OpenFile(path,flag|os.O_CREATE,mode);err != nil{
		_ = db.close()
		return nil,err
	}
	if err := 

}

