package storage

import("encoding/binary";"fmt")
func AppendRecord(p *Page,record []byte)(uint16,error){off:=HeaderSize;for off+2<=PageSize{n:=int(binary.LittleEndian.Uint16(p.Data[off:off+2]));if n==0{break};off+=2+n};if off+2+len(record)>PageSize{return 0,fmt.Errorf("page full")};binary.LittleEndian.PutUint16(p.Data[off:off+2],uint16(len(record)));copy(p.Data[off+2:],record);p.Dirty=true;return uint16(off+2),nil}
func Records(p *Page)[][]byte{var out [][]byte;for off:=HeaderSize;off+2<=PageSize;{n:=int(binary.LittleEndian.Uint16(p.Data[off:off+2]));if n==0||off+2+n>PageSize{break};out=append(out,append([]byte(nil),p.Data[off+2:off+2+n]...));off+=2+n};return out}
