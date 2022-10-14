//go:build windows

package mmap

import (
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

//mmap
func mmap(f uintptr, offset, len, prot, flags int64) ([]byte, uintptr, error) {
	cprot := syscall.PAGE_READONLY
	fileAccess := syscall.FILE_MAP_READ
	switch {
	case prot&COPY != 0:
		cprot = syscall.PAGE_WRITECOPY
		fileAccess = syscall.FILE_MAP_COPY
	case prot&RDWR != 0:
		cprot = syscall.PAGE_READWRITE
		fileAccess = syscall.FILE_MAP_WRITE
	}
	if prot&EXEC != 0 {
		cprot = cprot << 4
		fileAccess = fileAccess | syscall.FILE_MAP_EXECUTE
	}
	maxSize := offset + len
	maxSizeHight := uint32(maxSize >> 32)      //取高位32位
	maxSizeLow := uint32(maxSize & 0xFFFFFFFF) //取低位32位
	h, err := syscall.CreateFileMapping(syscall.Handle(f), nil, uint32(cprot), maxSizeHight, maxSizeLow, nil)
	if err != nil {
		return nil, 0, os.NewSyscallError("CreateFileMapping", err)
	}
	offsetHight := uint32(offset >> 32)      //取高位32位
	offsetLow := uint32(offset & 0xFFFFFFFF) //取低位32位

	addrPtr, err := syscall.MapViewOfFile(h, uint32(fileAccess), offsetHight, offsetLow, uintptr(len))
	if err != nil {
		syscall.CloseHandle(h)
		return nil, 0, os.NewSyscallError("MapViewOfFile", err)
	}
	//构建切片
	bb := make([]byte, 0)
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bb))
	slice.Len = int(len)
	slice.Cap = slice.Len
	slice.Data = addrPtr
	return bb, uintptr(h), nil
}

//flush
func flush(addrByte []byte) error {
	ptr, len, err := addrlen(addrByte)
	if err != nil {
		return err
	}
	err = syscall.FlushViewOfFile(ptr, len)
	return os.NewSyscallError("FlushViewOfFile", err)
}

//lock
func lock(addrByte []byte) error {
	ptr, len, err := addrlen(addrByte)
	if err != nil {
		return err
	}
	err = syscall.VirtualLock(ptr, len)
	return os.NewSyscallError("VirtualLock", err)
}

//unlock
func unlock(addrByte []byte) error {
	ptr, len, err := addrlen(addrByte)
	if err != nil {
		return err
	}
	err = syscall.VirtualUnlock(ptr, len)
	return os.NewSyscallError("VirtualUnlock", err)
}

//unmap
func unmap(addrByte []byte, handle uintptr) error {
	ptr, _, err := addrlen(addrByte)
	if err != nil {
		return err
	}
	err = syscall.UnmapViewOfFile(ptr)
	if err != nil {
		return os.NewSyscallError("UnmapViewOfFile", err)
	}
	if handle == 0 {
		return nil
	}
	err = syscall.CloseHandle(syscall.Handle(handle))
	return os.NewSyscallError("CloseHandle", err)
}
