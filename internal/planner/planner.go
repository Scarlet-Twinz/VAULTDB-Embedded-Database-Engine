package planner

type Kind string
const(TableScan Kind="table_scan";IndexScan Kind="index_scan";Filter Kind="filter";Projection Kind="projection";Join Kind="join";Aggregate Kind="aggregate";Sort Kind="sort";Limit Kind="limit")
type Plan struct{Kind Kind;Table string;Index string;Column string;Value string;Children []*Plan;Columns []string;GroupBy []string;OrderBy string;Limit int}
func Table(table string)*Plan{return &Plan{Kind:TableScan,Table:table}}
func Equality(table,index,column,value string)*Plan{return &Plan{Kind:IndexScan,Table:table,Index:index,Column:column,Value:value}}
func FilterPlan(child *Plan,column,value string)*Plan{return &Plan{Kind:Filter,Children:[]*Plan{child},Column:column,Value:value}}
func Project(child *Plan,cols []string)*Plan{return &Plan{Kind:Projection,Children:[]*Plan{child},Columns:cols}}
func LimitPlan(child *Plan,n int)*Plan{return &Plan{Kind:Limit,Children:[]*Plan{child},Limit:n}}
