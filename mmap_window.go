package mmap

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/windows"
)

//mmap
func mmap(f uintptr, offset, len, prot, flags int64) ([]byte, error) {
	cprot := windows.PAGE_READONLY
	fileAccess := windows.FILE_MAP_READ
	switch {
	case prot&COPY != 0:
		cprot = windows.PAGE_WRITECOPY
		fileAccess = windows.FILE_MAP_COPY
	case prot&RDWR != 0:
		cprot = windows.PAGE_READWRITE
		fileAccess = windows.FILE_MAP_WRITE
	}
	if prot&EXEC != 0 {
		cprot = cprot << 4
		fileAccess = fileAccess | windows.FILE_MAP_EXECUTE
	}
	maxSize := offset + len
	maxSizeHight := uint32(maxSize >> 32)      //取高位32位
	maxSizeLow := uint32(maxSize & 0xFFFFFFFF) //取低位32位
	h, err := windows.CreateFileMapping(windows.Handle(f), nil, uint32(cprot), maxSizeHight, maxSizeLow, nil)
	if err != nil {
		return nil, fmt.Errorf("CreateFileMapping %v", err)
	}
	offsetHight := uint32(offset >> 32)      //取高位32位
	offsetLow := uint32(offset & 0xFFFFFFFF) //取低位32位

	addrPtr, err := windows.MapViewOfFile(h, uint32(fileAccess), offsetHight, offsetLow, uintptr(len))
	if err != nil {
		windows.CloseHandle(h)
		return nil, fmt.Errorf("MapViewOfFile %v", err)
	}
	//构建切片
	bb := make([]byte, 0)
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bb))
	slice.Len = int(len)
	slice.Cap = slice.Len
	slice.Data = addrPtr
	return bb, nil
}
