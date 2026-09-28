package sql

import("fmt";"regexp";"strings")
type Statement interface{stmt()}
type CreateTable struct{Name string;Columns []Column};func(CreateTable)stmt(){}
type Column struct{Name,Type string;Nullable bool}
type Insert struct{Table string;Columns,Values []string};func(Insert)stmt(){}
type Select struct{Columns []string;Table string;WhereColumn,WhereValue string};func(Select)stmt(){}
type Delete struct{Table,WhereColumn,WhereValue string};func(Delete)stmt(){}
type Update struct{Table,SetColumn,SetValue,WhereColumn,WhereValue string};func(Update)stmt(){}
type Begin struct{};func(Begin)stmt(){}
type Commit struct{};func(Commit)stmt(){}
type Rollback struct{};func(Rollback)stmt(){}
type ShowTables struct{};func(ShowTables)stmt(){}
func clean(s string)string{return strings.TrimSpace(strings.TrimSuffix(s,";"))}
func Parse(input string)(Statement,error){
 s:=clean(input);u:=strings.ToUpper(s)
 switch{case u=="BEGIN":return Begin{},nil;case u=="COMMIT":return Commit{},nil;case u=="ROLLBACK":return Rollback{},nil;case u=="SHOW TABLES":return ShowTables{},nil}
 re:=regexp.MustCompile(`(?i)^CREATE\s+TABLE\s+([A-Za-z_][A-Za-z0-9_]*)\s*\((.+)\)$`)
 if m:=re.FindStringSubmatch(s);m!=nil{var cs []Column;for _,part:=range strings.Split(m[2],","){f:=strings.Fields(part);if len(f)<2{return nil,fmt.Errorf("invalid column definition %q",part)};cs=append(cs,Column{Name:f[0],Type:strings.ToUpper(f[1])})};return CreateTable{Name:m[1],Columns:cs},nil}
 re=regexp.MustCompile(`(?i)^INSERT\s+INTO\s+([A-Za-z_][A-Za-z0-9_]*)(?:\s*\(([^)]*)\))?\s+VALUES\s*\(([^)]*)\)$`)
 if m:=re.FindStringSubmatch(s);m!=nil{return Insert{Table:m[1],Columns:splitCSV(m[2]),Values:splitCSV(m[3])},nil}
 re=regexp.MustCompile(`(?i)^SELECT\s+(.+)\s+FROM\s+([A-Za-z_][A-Za-z0-9_]*)(?:\s+WHERE\s+([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+))?$`)
 if m:=re.FindStringSubmatch(s);m!=nil{return Select{Columns:splitCSV(m[1]),Table:m[2],WhereColumn:m[3],WhereValue:trimQuote(m[4])},nil}
 re=regexp.MustCompile(`(?i)^DELETE\s+FROM\s+([A-Za-z_][A-Za-z0-9_]*)(?:\s+WHERE\s+([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+))?$`)
 if m:=re.FindStringSubmatch(s);m!=nil{return Delete{Table:m[1],WhereColumn:m[2],WhereValue:trimQuote(m[3])},nil}
 re=regexp.MustCompile(`(?i)^UPDATE\s+([A-Za-z_][A-Za-z0-9_]*)\s+SET\s+([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+?)(?:\s+WHERE\s+([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.+))?$`)
 if m:=re.FindStringSubmatch(s);m!=nil{return Update{Table:m[1],SetColumn:m[2],SetValue:trimQuote(m[3]),WhereColumn:m[4],WhereValue:trimQuote(m[5])},nil}
 return nil,fmt.Errorf("unsupported SQL: %s",s)
}
func splitCSV(s string)[]string{if strings.TrimSpace(s)==""{return nil};var out []string;cur:="";quote:=byte(0);for i:=0;i<len(s);i++{c:=s[i];if c=='\''||c=='"'{if quote==0{quote=c}else if quote==c{quote=0}};if c==','&&quote==0{out=append(out,strings.TrimSpace(cur));cur=""}else{cur+=string(c)}};out=append(out,strings.TrimSpace(cur));return out}
func trimQuote(s string)string{s=strings.TrimSpace(s);if len(s)>=2&&((s[0]=='\''&&s[len(s)-1]=='\'')||(s[0]=='"'&&s[len(s)-1]=='"')){return s[1:len(s)-1]};return s}
