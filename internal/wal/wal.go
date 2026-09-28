package wal

import("bufio";"encoding/json";"fmt";"os";"sync";"time")
type Record struct{LSN uint64 `json:"lsn"`;TxID uint64 `json:"txid"`;Type string `json:"type"`;Data json.RawMessage `json:"data,omitempty"`;Time time.Time `json:"time"`}
type WAL struct{mu sync.Mutex;f *os.File;next uint64}
func Open(path string)(*WAL,error){f,e:=os.OpenFile(path,os.O_CREATE|os.O_RDWR|os.O_APPEND,0600);if e!=nil{return nil,e};return &WAL{f:f,next:uint64(time.Now().UnixNano())},nil}
func(w *WAL)Append(tx uint64,typ string,data any)(uint64,error){w.mu.Lock();defer w.mu.Unlock();b,e:=json.Marshal(data);if e!=nil{return 0,e};r:=Record{LSN:w.next,TxID:tx,Type:typ,Data:b,Time:time.Now().UTC()};w.next++;line,e:=json.Marshal(r);if e!=nil{return 0,e};if _,e=w.f.Write(append(line,byte(10)));e!=nil{return 0,e};if e=w.f.Sync();e!=nil{return 0,e};return r.LSN,nil}
func(w *WAL)Scan()([]Record,error){w.mu.Lock();defer w.mu.Unlock();if _,e:=w.f.Seek(0,0);e!=nil{return nil,e};s:=bufio.NewScanner(w.f);var out []Record;for s.Scan(){var r Record;if e:=json.Unmarshal(s.Bytes(),&r);e!=nil{return nil,fmt.Errorf("invalid WAL record: %w",e)};out=append(out,r)};_,_=w.f.Seek(0,2);return out,s.Err()}
func(w *WAL)Close()error{return w.f.Close()}
