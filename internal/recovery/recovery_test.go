package recovery
import("testing";"github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/wal")
func TestCommitted(t *testing.T){rs:=[]wal.Record{{TxID:1,Type:"BEGIN"},{TxID:1,Type:"INSERT"},{TxID:1,Type:"COMMIT"},{TxID:2,Type:"INSERT"},{TxID:2,Type:"ROLLBACK"}};out:=Committed(rs);if len(out)!=1||out[0].TxID!=1{t.Fatalf("unexpected %v",out)}}
