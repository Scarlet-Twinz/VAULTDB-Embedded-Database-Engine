package storage

import (
	"encoding/binary"
	"fmt"
)

const PageSize = 4096
const HeaderSize = 16

type PageID uint64

type Page struct {
	ID    PageID
	Data  [PageSize]byte
	Dirty bool
}

func NewPage(id PageID) *Page { return &Page{ID: id} }
func (p *Page) EncodeHeader() {
	binary.LittleEndian.PutUint64(p.Data[0:8], uint64(p.ID))
	binary.LittleEndian.PutUint32(p.Data[8:12], uint32(PageSize))
}
func DecodePage(buf []byte) (*Page, error) {
	if len(buf) != PageSize {
		return nil, fmt.Errorf("invalid page size: %d", len(buf))
	}
	var p Page
	p.ID = PageID(binary.LittleEndian.Uint64(buf))
	copy(p.Data[:], buf)
	return &p, nil
}
