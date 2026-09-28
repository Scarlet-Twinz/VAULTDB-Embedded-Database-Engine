package main

import("bufio";"flag";"fmt";"os";"strings";"github.com/Scarlet-Twinz/VAULTDB-Embedded-Database-Engine/internal/engine")
func main(){dir:=flag.String("database","./vaultdb-data","database directory");exec:=flag.String("execute","","execute one SQL statement");flag.Parse();db,e:=engine.Open(*dir);if e!=nil{fmt.Fprintln(os.Stderr,e);return};defer db.Close();if *exec!=""{run(db,*exec);return};fmt.Println("VAULTDB — embedded relational database");fmt.Println("Type .help for commands.");s:=bufio.NewScanner(os.Stdin);for{fmt.Print("vaultdb> ");if !s.Scan(){break};q:=strings.TrimSpace(s.Text());if q==".exit"{break};if q==".help"{fmt.Println("SQL: CREATE TABLE, INSERT, SELECT, UPDATE, DELETE, BEGIN, COMMIT, ROLLBACK, SHOW TABLES");continue};run(db,q)}}
func run(db *engine.Engine,q string){rows,e:=db.Exec(q);if e!=nil{fmt.Println("ERROR:",e);return};for _,r:=range rows{fmt.Println(r)}}
