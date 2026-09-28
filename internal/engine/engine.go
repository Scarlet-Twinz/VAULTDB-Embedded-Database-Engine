package engine

import("encoding/json";"fmt";"os";"path/filepath";"strconv";"strings";"sync"
 "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/catalog"
 "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/sql"
 "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/storage"
 "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/transaction"
 "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/wal")
type Row map[string]string
type Engine struct{mu sync.Mutex;dir string;file *storage.File;buf *storage.BufferPool;catalog *catalog.Catalog;wal *wal.WAL;tx *transaction.Manager;active *transaction.Tx}
func Open(dir string)(*Engine,error){if e:=os.MkdirAll(dir,0700);e!=nil{return nil,e};f,e:=storage.Open(filepath.Join(dir,"data.db"));if e!=nil{return nil,e};w,e:=wal.Open(filepath.Join(dir,"wal.log"));if e!=nil{f.Close();return nil,e};c,e:=catalog.Open(filepath.Join(dir,"catalog.json"));if e!=nil{f.Close();w.Close();return nil,e};return &Engine{dir:dir,file:f,buf:storage.NewBufferPool(f,64),catalog:c,wal:w,tx:transaction.NewManager(w)},nil}
func(e *Engine)Close()error{e.mu.Lock();defer e.mu.Unlock();if x:=e.buf.FlushAll();x!=nil{return x};if x:=e.catalog.Save();x!=nil{return x};if x:=e.wal.Close();x!=nil{return x};return e.file.Close()}
func(e *Engine)Exec(q string)([]Row,error){e.mu.Lock();defer e.mu.Unlock();s,err:=sql.Parse(q);if err!=nil{return nil,err};switch x:=s.(type){
case sql.Begin:t,er:=e.tx.Begin();if er==nil{e.active=t};return nil,er
case sql.Commit:if e.active==nil{return nil,fmt.Errorf("no active transaction")};er:=e.tx.Commit(e.active);if er==nil{e.active=nil};return nil,er
case sql.Rollback:if e.active==nil{return nil,fmt.Errorf("no active transaction")};er:=e.tx.Rollback(e.active);if er==nil{e.active=nil};return nil,er
case sql.ShowTables:n:=e.catalog.ListTables();r:=make([]Row,0,len(n));for _,v:=range n{r=append(r,Row{"table":v})};return r,nil
case sql.CreateTable:return nil,e.createTable(x)
case sql.Insert:return nil,e.insert(x)
case sql.Select:return e.selectRows(x)
case sql.Update:return nil,e.update(x)
case sql.Delete:return nil,e.delete(x)}
return nil,fmt.Errorf("unsupported statement")}
func(e *Engine)createTable(x sql.CreateTable)error{cols:=make([]catalog.Column,len(x.Columns));for i,c:=range x.Columns{cols[i]=catalog.Column{Name:c.Name,Type:c.Type,Nullable:c.Nullable}};if er:=e.catalog.CreateTable(catalog.Table{Name:x.Name,Columns:cols});er!=nil{return er};if e.active!=nil{_,_=e.wal.Append(e.active.ID,"CREATE_TABLE",x)};return e.catalog.Save()}
func(e *Engine)insert(x sql.Insert)error{t,ok:=e.catalog.GetTable(x.Table);if !ok{return fmt.Errorf("table %q not found",x.Table)};cols:=x.Columns;if len(cols)==0{for _,c:=range t.Columns{cols=append(cols,c.Name)}};if len(cols)!=len(x.Values){return fmt.Errorf("column/value count mismatch")};r:=Row{};for i,c:=range cols{r[c]=trim(x.Values[i])};if er:=e.appendRow(&t,r);er!=nil{return er};if e.active!=nil{_,_=e.wal.Append(e.active.ID,"INSERT",r)};return nil}
func(e *Engine)appendRow(t *catalog.Table,r Row)error{b,er:=json.Marshal(r);if er!=nil{return er};if len(t.Pages)==0{p:=storage.NewPage(0);if _,er=storage.AppendRecord(p,b);er!=nil{return er};if er=e.file.WritePage(p);er!=nil{return er};t.Pages=[]uint64{0};e.catalog.TablesUnsafeSet(*t);return e.catalog.Save()}
pid:=t.Pages[len(t.Pages)-1];p,er:=e.buf.Fetch(storage.PageID(pid));if er!=nil{return er}
if _,er=storage.AppendRecord(p,b);er==nil{e.buf.Unpin(storage.PageID(pid),true);e.catalog.TablesUnsafeSet(*t);return e.catalog.Save()}
e.buf.Unpin(storage.PageID(pid),false);newID:=storage.PageID(len(t.Pages));p=storage.NewPage(newID);if _,er=storage.AppendRecord(p,b);er!=nil{return er};if er=e.file.WritePage(p);er!=nil{return er};t.Pages=append(t.Pages,uint64(newID));e.catalog.TablesUnsafeSet(*t);return e.catalog.Save()}
func(e *Engine)selectRows(x sql.Select)([]Row,error){t,ok:=e.catalog.GetTable(x.Table);if !ok{return nil,fmt.Errorf("table %q not found",x.Table)};var out []Row;for _,pid:=range t.Pages{p,er:=e.buf.Fetch(storage.PageID(pid));if er!=nil{return nil,er};for _,b:=range storage.Records(p){var r Row;if json.Unmarshal(b,&r)==nil&&match(r,x.WhereColumn,x.WhereValue){if len(x.Columns)==1&&x.Columns[0]=="*"{out=append(out,r)}else{nr:=Row{};for _,c:=range x.Columns{nr[c]=r[c]};out=append(out,nr)}}};e.buf.Unpin(storage.PageID(pid),false)};return out,nil}
func(e *Engine)update(x sql.Update)error{if _,er:=e.selectRows(sql.Select{Columns:[]string{"*"},Table:x.Table,WhereColumn:x.WhereColumn,WhereValue:x.WhereValue});er!=nil{return er};return e.rewriteTable(x.Table,func(r Row)Row{if match(r,x.WhereColumn,x.WhereValue){r[x.SetColumn]=trim(x.SetValue)};return r})}
func(e *Engine)delete(x sql.Delete)error{return e.rewriteTable(x.Table,func(r Row)Row{if match(r,x.WhereColumn,x.WhereValue){return nil};return r})}
func(e *Engine)rewriteTable(name string,fn func(Row)Row)error{t,ok:=e.catalog.GetTable(name);if !ok{return fmt.Errorf("table %q not found",name)};var rows []Row;for _,pid:=range t.Pages{p,er:=e.buf.Fetch(storage.PageID(pid));if er!=nil{return er};for _,b:=range storage.Records(p){var r Row;if json.Unmarshal(b,&r)==nil{if r=fn(r);r!=nil{rows=append(rows,r)}}};e.buf.Unpin(storage.PageID(pid),false)};t.Pages=nil;for _,r:=range rows{if er:=e.appendRow(&t,r);er!=nil{return er}};e.catalog.TablesUnsafeSet(t);return e.catalog.Save()}
func match(r Row,c,v string)bool{if c==""{return true};return r[c]==v}
func trim(s string)string{s=strings.TrimSpace(s);if len(s)>=2&&((s[0]=='\''&&s[len(s)-1]=='\'')||(s[0]=='"'&&s[len(s)-1]=='"')){return s[1:len(s)-1]};if i,e:=strconv.Atoi(s);e==nil{return strconv.Itoa(i)};return s}
