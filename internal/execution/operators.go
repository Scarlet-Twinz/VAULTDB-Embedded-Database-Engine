package execution

import("fmt";"sort";"strconv")
type Row map[string]string
type Operator interface{Next()([]Row,error)}
type Slice struct{Rows []Row;done bool}
func(s *Slice)Next()([]Row,error){if s.done{return nil,nil};s.done=true;return s.Rows,nil}
type FilterOp struct{Child Operator;Column,Value string}
func(o *FilterOp)Next()([]Row,error){rs,e:=o.Child.Next();if e!=nil{return nil,e};out:=make([]Row,0,len(rs));for _,r:=range rs{if o.Column==""||r[o.Column]==o.Value{out=append(out,r)}};return out,nil}
type ProjectOp struct{Child Operator;Columns []string}
func(o *ProjectOp)Next()([]Row,error){rs,e:=o.Child.Next();if e!=nil{return nil,e};if len(o.Columns)==0||len(o.Columns)==1&&o.Columns[0]=="*"{return rs,nil};out:=make([]Row,0,len(rs));for _,r:=range rs{n:=Row{};for _,c:=range o.Columns{n[c]=r[c]};out=append(out,n)};return out,nil}
type LimitOp struct{Child Operator;Limit int}
func(o *LimitOp)Next()([]Row,error){rs,e:=o.Child.Next();if e!=nil{return nil,e};if o.Limit>=0&&len(rs)>o.Limit{return rs[:o.Limit],nil};return rs,nil}
func Join(left,right []Row,leftKey,rightKey string)[]Row{out:=[]Row{};for _,l:=range left{for _,r:=range right{if l[leftKey]==r[rightKey]{n:=Row{};for k,v:=range l{n[k]=v};for k,v:=range r{if _,ok:=n[k];ok{n["right."+k]=v}else{n[k]=v}};out=append(out,n)}}};return out}
func Group(rows []Row,keys []string,measure string)[]Row{groups:=map[string][]Row{};order:=[]string{};for _,r:=range rows{k:="";for _,c:=range keys{k+="|"+r[c]};if _,ok:=groups[k];!ok{order=append(order,k)};groups[k]=append(groups[k],r)};out:=make([]Row,0,len(order));for _,k:=range order{rs:=groups[k];n:=Row{};for _,c:=range keys{n[c]=rs[0][c]};parts:=splitMeasure(measure);if len(parts)==2{fn,col:=parts[0],parts[1];n[measure]=aggregate(fn,col,rs)};out=append(out,n)};return out}
func aggregate(fn,col string,rows []Row)string{if fn=="COUNT"{return strconv.Itoa(len(rows))};sum:=0.0;min,max:="","";for _,r:=range rows{v:=r[col];x,e:=strconv.ParseFloat(v,64);if e!=nil{continue};sum+=x;if min==""||x<parse(min){min=v};if max==""||x>parse(max){max=v}};switch fn{case"SUM":return fmt.Sprintf("%g",sum);case"AVG":if len(rows)==0{return"0"};return fmt.Sprintf("%g",sum/float64(len(rows)));case"MIN":return min;case"MAX":return max};return""}
func parse(s string)float64{x,_:=strconv.ParseFloat(s,64);return x}
func splitMeasure(s string)[]string{for _,fn:=range []string{"COUNT","SUM","AVG","MIN","MAX"}{p:=fn+"(";if len(s)>len(p)&&len(s)>0&&len(s)>=len(p)+1&&s[:len(p)]==p&&s[len(s)-1]==')'{return []string{fn,s[len(p):len(s)-1]}}};return nil}
func Sort(rows []Row,column string,desc bool)[]Row{out:=append([]Row(nil),rows...);sort.SliceStable(out,func(i,j int)bool{if desc{return out[i][column]>out[j][column]};return out[i][column]<out[j][column]});return out}
