//go:build darwin || dragonfly || freebsd || linux || openbsd || solaris || netbsd

package mmap

import (
	"golang.org/x/sys/unix"
)

//mmap
func mmap(f uintptr, offset, len, inprot, inflags int64) ([]byte, uintptr, error) {
	flags := unix.MAP_SHARED
	prot := unix.PROT_READ
	switch {
	case inprot&COPY != 0:
		prot |= unix.PROT_WRITE
		flags = unix.MAP_PRIVATE
	case inprot&RDWR != 0:
		prot |= unix.PROT_WRITE
	}
	if inprot&EXEC != 0 {
		prot |= unix.PROT_EXEC
	}
	if inflags&ANON != 0 {
		flags |= unix.MAP_ANON
	}

	b, err := unix.Mmap(int(f), offset, int(len), prot, flags)
	if err != nil {
		return nil, 0, err
	}
	return b, 0, nil
}

//flush
func flush(addrByte []byte) error {
	return unix.Msync([]byte(addrByte), unix.MS_SYNC)
}

//lock
func lock(addrByte []byte) error {
	return unix.Mlock([]byte(addrByte))
}

//unlock
func unlock(addrByte []byte) error {
	return unix.Munlock([]byte(addrByte))
}

//unmap
func unmap(addrByte []byte, handle uintptr) error {
	return unix.Munmap([]byte(addrByte))
}
