package mmap

import (
	"errors"
	"reflect"
	"unsafe"
)

//base github.com/edsrzf/mmap-go
const (
	// RDONLY maps the memory read-only.
	// Attempts to write to the MMap object will result in undefined behavior.
	RDONLY = 0
	// RDWR maps the memory as read-write. Writes to the MMap object will update the
	// underlying file.
	RDWR = 1 << iota
	// COPY maps the memory as copy-on-write. Writes to the MMap object will affect
	// memory, but the underlying file will remain unchanged.
	COPY
	// If EXEC is set, the mapped memory is marked as executable.
	EXEC
)

const (
	// If the ANON flag is set, the mapped memory will not be backed by a file.
	ANON = 1 << iota
)

func addrlen(addrByte []byte) (ptr uintptr, len uintptr, err error) {
	if addrByte == nil {
		ptr, len, err = 0, 0, errors.New("addrByte nil")
		return
	}
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&addrByte))
	ptr, len, err = slice.Data, uintptr(slice.Len), nil
	return
}
