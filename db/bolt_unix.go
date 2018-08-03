package db

import (
	"os"
	"time"
	"syscall"
)
//acquires an advisory lock on a file descriptor.
func flock(db *DB,mode os.FileMode,exclusive bool,timeout time.Duration)error{
	var t time.Time
	for{
		if t.IsZero() {
			t = time.Now()
		}else if timeout > 0 && time.Since(t) > timeout {
			return ErrTimeout
		}
		flag := syscall.LOCK_SH
		if exclusive {
			flag = syscall.LOCK_EX
		}
	}
}
