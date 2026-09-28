package index
import "testing"
func TestBTreeInsertSearchDelete(t *testing.T){b:=New(4);for i,k:=range []string{"d","a","c","b","e","f"}{if e:=b.Insert(k,uint64(i));e!=nil{t.Fatal(e)}};if v,ok:=b.Search("c");!ok||v!=2{t.Fatalf("search %v %v",v,ok)};if !b.Delete("c"){t.Fatal("delete failed")};if _,ok:=b.Search("c");ok{t.Fatal("deleted key found")}}
