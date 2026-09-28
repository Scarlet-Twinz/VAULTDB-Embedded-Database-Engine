package concurrency
import "testing"
func TestLockConflict(t *testing.T){l:=NewLockManager();if e:=l.Acquire("t",Write);e!=nil{t.Fatal(e)};if e:=l.Acquire("t",Read);e==nil{t.Fatal("expected conflict")};l.Release("t",Write);if e:=l.Acquire("t",Read);e!=nil{t.Fatal(e)};l.Release("t",Read)}
