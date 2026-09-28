package metrics
import("testing";"time")
func TestMetrics(t *testing.T){var m Metrics;m.IncPagesRead();m.IncBufferHit();m.AddQuery(time.Millisecond);s:=m.Snapshot();if s.PagesRead!=1||s.BufferHits!=1||s.Queries!=1{t.Fatal(s)}}
