package db

import "errors"

var ErrTimeout = errors.New("timeout")
var ErrInvalid = errors.New("invalid database")
var ErrVersionMismatch = errors.New("version mismatch")
var ErrChecksum = errors.New("checksum error")
