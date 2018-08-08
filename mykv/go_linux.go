package mykv

import "syscall"

func fdatasync (db *DB)error{
	return syscall.Fdatasync(int(db.file.Fd()))
}
