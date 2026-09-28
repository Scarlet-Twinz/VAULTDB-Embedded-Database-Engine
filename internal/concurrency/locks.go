package concurrency

import("fmt";"sync")
type Mode uint8
const(Read Mode=1;Write Mode=2)
type LockManager struct{mu sync.Mutex;locks map[string]*entry}
type entry struct{readers int;writer bool}
func NewLockManager()*LockManager{return &LockManager{locks:map[string]*entry{}}}
func(l *LockManager)Acquire(key string,mode Mode)error{l.mu.Lock();defer l.mu.Unlock();e:=l.locks[key];if e==nil{e=&entry{};l.locks[key]=e};if mode==Read{if e.writer{return fmt.Errorf("write lock held on %s",key)};e.readers++;return nil};if e.writer||e.readers>0{return fmt.Errorf("lock conflict on %s",key)};e.writer=true;return nil}
func(l *LockManager)Release(key string,mode Mode){l.mu.Lock();defer l.mu.Unlock();if e:=l.locks[key];e!=nil{if mode==Read&&e.readers>0{e.readers--};if mode==Write{e.writer=false};if e.readers==0&&!e.writer{delete(l.locks,key)}}}
