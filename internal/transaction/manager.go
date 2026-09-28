package transaction

import("fmt";"sync";"github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/wal")
type State string
const(Begin State="BEGIN";Committed State="COMMITTED";Aborted State="ABORTED")
type Tx struct{ID uint64;State State}
type Manager struct{mu sync.Mutex;next uint64;active map[uint64]*Tx;wal *wal.WAL}
func NewManager(w *wal.WAL)*Manager{return &Manager{next:1,active:map[uint64]*Tx{},wal:w}}
func(m *Manager)Begin()(*Tx,error){m.mu.Lock();defer m.mu.Unlock();id:=m.next;m.next++;t:=&Tx{ID:id,State:Begin};m.active[id]=t;if _,e:=m.wal.Append(id,"BEGIN",nil);e!=nil{return nil,e};return t,nil}
func(m *Manager)Commit(t *Tx)error{m.mu.Lock();defer m.mu.Unlock();if t.State!=Begin{return fmt.Errorf("transaction %d is not active",t.ID)};if _,e:=m.wal.Append(t.ID,"COMMIT",nil);e!=nil{return e};t.State=Committed;delete(m.active,t.ID);return nil}
func(m *Manager)Rollback(t *Tx)error{m.mu.Lock();defer m.mu.Unlock();if t.State!=Begin{return fmt.Errorf("transaction %d is not active",t.ID)};if _,e:=m.wal.Append(t.ID,"ROLLBACK",nil);e!=nil{return e};t.State=Aborted;delete(m.active,t.ID);return nil}
