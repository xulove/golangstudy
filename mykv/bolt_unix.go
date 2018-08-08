package mykv

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

//acquires an advisory lock on a file descriptor.
func flock(db *DB, mode os.FileMode, exclusive bool, timeout time.Duration) error {
	var t time.Time
	for {
		if t.IsZero() {
			t = time.Now()
		} else if timeout > 0 && time.Since(t) > timeout {
			return ErrTimeout
		}
		flag := syscall.LOCK_SH
		if exclusive {
			flag = syscall.LOCK_EX
		}
		err := syscall.Flock(int(db.file.Fd()), flag|syscall.LOCK_NB)
		if err == nil {
			return nil
		} else if err != syscall.EWOULDBLOCK {
			return err
		}
		// 等待一段时间重试
		time.Sleep(50 * time.Millisecond)

	}
}

// 内存映射
func mmap(db *DB, sz int) error {
	// Map the data file to memory.
	b, err := syscall.Mmap(int(db.file.Fd()), 0, sz, syscall.PROT_READ, syscall.MAP_SHARED|db.MmapFlags)
	if err != nil {
		return err
	}
	//建议内核随机访问mmap
	if err := madvise(b, syscall.MADV_RANDOM); err != nil {
		return fmt.Errorf("madvise: %s", err)
	}
	db.dataref = b
	db.data = (*[maxMapSize]byte)(unsafe.Pointer(&b[0]))
	db.datasz = sz
	return nil
}

// NOTE: This function is copied from stdlib because it is not available on darwin.
func madvise(b []byte, advice int) (err error) {
	_, _, e1 := syscall.Syscall(syscall.SYS_MADVISE, uintptr(unsafe.Pointer(&b[0])), uintptr(len(b)), uintptr(advice))
	if e1 != 0 {
		err = e1
	}
	return
}

// 清除内存中的映射
func munmap(db *DB) error {
	if db.dataref == nil {
		return nil
	}
	// Unmap using the original byte slice.
	err := syscall.Munmap(db.dataref)
	db.dataref = nil
	db.data = nil
	db.datasz = 0
	return err
}
