package mykv

type bucket struct {
	root     pgid
	sequence uint64
}
