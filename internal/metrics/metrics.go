package metrics
import("sync";"time")
type Snapshot struct{PagesRead,PagesWritten,BufferHits,BufferMisses,IndexLookups,TableScans,Transactions uint64;Queries uint64;QueryNanos int64}
type Metrics struct{mu sync.RWMutex;s Snapshot}
func(m *Metrics)AddQuery(d time.Duration){m.mu.Lock();m.s.Queries++;m.s.QueryNanos+=d.Nanoseconds();m.mu.Unlock()}
func(m *Metrics)IncPagesRead(){m.mu.Lock();m.s.PagesRead++;m.mu.Unlock()}
func(m *Metrics)IncPagesWritten(){m.mu.Lock();m.s.PagesWritten++;m.mu.Unlock()}
func(m *Metrics)IncBufferHit(){m.mu.Lock();m.s.BufferHits++;m.mu.Unlock()}
func(m *Metrics)IncBufferMiss(){m.mu.Lock();m.s.BufferMisses++;m.mu.Unlock()}
func(m *Metrics)IncIndexLookup(){m.mu.Lock();m.s.IndexLookups++;m.mu.Unlock()}
func(m *Metrics)IncTableScan(){m.mu.Lock();m.s.TableScans++;m.mu.Unlock()}
func(m *Metrics)IncTransaction(){m.mu.Lock();m.s.Transactions++;m.mu.Unlock()}
func(m *Metrics)Snapshot()Snapshot{m.mu.RLock();defer m.mu.RUnlock();return m.s}
