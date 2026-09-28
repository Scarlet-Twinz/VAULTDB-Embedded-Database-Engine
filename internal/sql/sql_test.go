package sql
import "testing"
func TestParseCreateInsertSelect(t *testing.T){tests:=[]string{"CREATE TABLE users (id INT, name TEXT)","INSERT INTO users VALUES (1, 'Ada')","SELECT * FROM users WHERE id = 1"};for _,q:=range tests{if _,e:=Parse(q);e!=nil{t.Fatalf("%s: %v",q,e)}}}
