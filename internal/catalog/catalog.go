package catalog

import("encoding/json";"fmt";"os";"sync")
type Column struct{Name string `json:"name"`;Type string `json:"type"`;Nullable bool `json:"nullable"`}
type Table struct{Name string `json:"name"`;Columns []Column `json:"columns"`;Pages []uint64 `json:"pages"`}
type Catalog struct{mu sync.RWMutex;Path string `json:"-"`;Tables map[string]Table `json:"tables"`}
func Open(path string)(*Catalog,error){c:=&Catalog{Path:path,Tables:map[string]Table{}};b,e:=os.ReadFile(path);if os.IsNotExist(e){return c,nil};if e!=nil{return nil,e};if e=json.Unmarshal(b,c);e!=nil{return nil,e};if c.Tables==nil{c.Tables=map[string]Table{}};return c,nil}
func(c *Catalog)Save()error{c.mu.RLock();defer c.mu.RUnlock();b,e:=json.MarshalIndent(c,"","  ");if e!=nil{return e};tmp:=c.Path+".tmp";if e=os.WriteFile(tmp,b,0600);e!=nil{return e};return os.Rename(tmp,c.Path)}
func(c *Catalog)CreateTable(t Table)error{c.mu.Lock();defer c.mu.Unlock();if _,ok:=c.Tables[t.Name];ok{return fmt.Errorf("table %q already exists",t.Name)};c.Tables[t.Name]=t;return nil}
func(c *Catalog)GetTable(name string)(Table,bool){c.mu.RLock();defer c.mu.RUnlock();t,ok:=c.Tables[name];return t,ok}
func(c *Catalog)ListTables()[]string{c.mu.RLock();defer c.mu.RUnlock();o:=make([]string,0,len(c.Tables));for n:=range c.Tables{o=append(o,n)};return o}
func(c *Catalog)TablesUnsafeSet(t Table){c.mu.Lock();defer c.mu.Unlock();c.Tables[t.Name]=t}
