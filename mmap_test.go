package mmap

import (
	"fmt"
	"os"
	"testing"
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
	addr, err := mmap(f.Fd(), 0, 1<<30, RDWR, 0)
	if err != nil {
		panic(err)
	}
	c := []byte("hello world")
	copy(addr[:len(c)], c)
	fmt.Println(string(addr[:len(c)]))
}
