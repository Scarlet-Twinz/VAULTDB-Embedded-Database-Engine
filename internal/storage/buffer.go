package storage

import("container/list";"fmt";"sync")
type frame struct{p *Page;pin int;elem *list.Element}
type BufferPool struct{mu sync.Mutex;file *File;cap int;frames map[PageID]*frame;lru *list.List;Hits,Misses uint64}
func NewBufferPool(file *File,cap int)*BufferPool{if cap<1{cap=1};return &BufferPool{file:file,cap:cap,frames:map[PageID]*frame{},lru:list.New()}}
func(b *BufferPool)Fetch(id PageID)(*Page,error){b.mu.Lock();defer b.mu.Unlock();if fr:=b.frames[id];fr!=nil{b.Hits++;b.lru.MoveToFront(fr.elem);fr.pin++;return fr.p,nil};b.Misses++;p,e:=b.file.ReadPage(id);if e!=nil{return nil,e};if len(b.frames)>=b.cap{var victim *frame;for e:=b.lru.Back();e!=nil;e=e.Prev(){v:=e.Value.(*frame);if v.pin==0{victim=v;break}};if victim==nil{return nil,fmt.Errorf("buffer pool exhausted")};if victim.p.Dirty{if e:=b.file.WritePage(victim.p);e!=nil{return nil,e}};delete(b.frames,victim.p.ID);b.lru.Remove(victim.elem)};fr:=&frame{p:p,pin:1};fr.elem=b.lru.PushFront(fr);b.frames[id]=fr;return p,nil}
func(b *BufferPool)Unpin(id PageID,dirty bool){b.mu.Lock();defer b.mu.Unlock();if fr:=b.frames[id];fr!=nil{if dirty{fr.p.Dirty=true};if fr.pin>0{fr.pin--}}}
func(b *BufferPool)FlushAll()error{b.mu.Lock();defer b.mu.Unlock();for _,fr:=range b.frames{if fr.p.Dirty{if e:=b.file.WritePage(fr.p);e!=nil{return e};fr.p.Dirty=false}};return nil}
