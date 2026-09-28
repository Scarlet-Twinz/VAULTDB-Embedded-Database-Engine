package storage

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type File struct {
	mu sync.Mutex
	f  *os.File
}

func Open(path string) (*File, error) {
	f, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if e != nil { return nil, e }
	return &File{f: f}, nil
}
func (f *File) Close() error { return f.f.Close() }
func (f *File) ReadPage(id PageID) (*Page, error) {
	f.mu.Lock(); defer f.mu.Unlock()
	buf := make([]byte, PageSize)
	_, e := f.f.ReadAt(buf, int64(id)*PageSize)
	if e != nil {
		if e == io.EOF { return nil, io.ErrUnexpectedEOF }
		return nil, e
	}
	return DecodePage(buf)
}
func (f *File) WritePage(p *Page) error {
	f.mu.Lock(); defer f.mu.Unlock()
	p.EncodeHeader()
	n, e := f.f.WriteAt(p.Data[:], int64(p.ID)*PageSize)
	if e == nil && n != PageSize { return fmt.Errorf("short page write: %d", n) }
	if e == nil { e = f.f.Sync() }
	return e
}
func (f *File) PageCount() (uint64, error) {
	s, e := f.f.Stat(); if e != nil { return 0, e }
	return uint64((s.Size()+PageSize-1)/PageSize), nil
}
