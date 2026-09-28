package recovery

import "github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/wal"
type TxState struct{Committed bool;Aborted bool}
func Analyze(records []wal.Record)map[uint64]TxState{out:=map[uint64]TxState{};for _,r:=range records{s:=out[r.TxID];switch r.Type{case "COMMIT":s.Committed=true;case "ROLLBACK":s.Aborted=true};out[r.TxID]=s};return out}
func Committed(records []wal.Record)[]wal.Record{state:=Analyze(records);var out []wal.Record;for _,r:=range records{if state[r.TxID].Committed&&r.Type!="COMMIT"&&r.Type!="BEGIN"{out=append(out,r)}};return out}
