package mmap

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestIntMmap(t *testing.T) {
	f, err := os.OpenFile("./test.mempry", os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		panic(err)
	}

	defer f.Close()
	//windows.MapViewOfFile(windows.Handle(f.Fd()), 0)
	// syscall.CreateFileMapping()
	// syscall.MapViewOfFile()
	// windows.CreateFileMapping()
	// windows.PAGE_READONLY
	addr, h, err := mmap(f.Fd(), 0, 100<<30, RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer unmap(addr, h)
	c := []byte("abc你好世界！世界和平oopp")
	copy(addr[:len(c)], c)
	fmt.Println(string(addr[:len(c)]))
	flush(addr)
	time.Sleep(time.Second * 2)
}
