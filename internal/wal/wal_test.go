package wal
import "testing"
func TestWAL(t *testing.T){w,e:=Open(t.TempDir()+"/wal.log");if e!=nil{t.Fatal(e)};defer w.Close();if _,e=w.Append(1,"BEGIN",nil);e!=nil{t.Fatal(e)};rs,e:=w.Scan();if e!=nil||len(rs)!=1{t.Fatalf("scan: %v %d",e,len(rs))}}
