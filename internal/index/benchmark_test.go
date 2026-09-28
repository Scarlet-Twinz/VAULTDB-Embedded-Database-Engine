package index
import("strconv";"testing")
func BenchmarkBTreeInsert(b *testing.B){for n:=0;n<b.N;n++{t:=New(32);for i:=0;i<100;i++{_ = t.Insert(strconv.Itoa(i),uint64(i))}}}
