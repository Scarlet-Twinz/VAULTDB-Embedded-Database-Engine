package index

import (
 "encoding/json"
 "fmt"
 "os"
 "sort"
 "sync"
)

type Entry struct{Key string `json:"key"`;Value uint64 `json:"value"`}
type node struct{Leaf bool `json:"leaf"`;Keys []string `json:"keys"`;Values []uint64 `json:"values"`;Children []*node `json:"children"`;Next *node `json:"-"`}
type BTree struct{mu sync.RWMutex;Order int;Root *node;Path string}
func New(order int)*BTree{if order<3{order=3};return &BTree{Order:order,Root:&node{Leaf:true}}}
func Open(path string,order int)(*BTree,error){b:=New(order);b.Path=path;raw,e:=os.ReadFile(path);if os.IsNotExist(e){return b,nil};if e!=nil{return nil,e};if e=json.Unmarshal(raw,b);e!=nil{return nil,e};return b,nil}
func(b *BTree)Save()error{b.mu.RLock();defer b.mu.RUnlock();if b.Path==""{return fmt.Errorf("index path is empty")};raw,e:=json.Marshal(b);if e!=nil{return e};tmp:=b.Path+".tmp";if e=os.WriteFile(tmp,raw,0600);e!=nil{return e};return os.Rename(tmp,b.Path)}
func(b *BTree)Search(key string)(uint64,bool){b.mu.RLock();defer b.mu.RUnlock();n:=b.Root;for !n.Leaf{i:=sort.SearchStrings(n.Keys,key);if i<len(n.Keys)&&key>=n.Keys[i]{i++};n=n.Children[i]};i:=sort.SearchStrings(n.Keys,key);if i<len(n.Keys)&&n.Keys[i]==key{return n.Values[i],true};return 0,false}
func(b *BTree)Insert(key string,value uint64)error{b.mu.Lock();defer b.mu.Unlock();if b.Root==nil{b.Root=&node{Leaf:true}};if len(b.Root.Keys)>=b.Order-1{old:=b.Root;b.Root=&node{Children:[]*node{old}};b.splitChild(b.Root,0)};return b.insertNonFull(b.Root,key,value)}
func(b *BTree)insertNonFull(n *node,key string,value uint64)error{if n.Leaf{i:=sort.SearchStrings(n.Keys,key);if i<len(n.Keys)&&n.Keys[i]==key{return fmt.Errorf("duplicate index key %q",key)};n.Keys=append(n.Keys,"");copy(n.Keys[i+1:],n.Keys[i:]);n.Keys[i]=key;n.Values=append(n.Values,0);copy(n.Values[i+1:],n.Values[i:]);n.Values[i]=value;return nil};i:=sort.SearchStrings(n.Keys,key);if len(n.Children[i].Keys)>=b.Order-1{b.splitChild(n,i);if key>n.Keys[i]{i++}else if key==n.Keys[i]{return fmt.Errorf("duplicate index key %q",key)}};return b.insertNonFull(n.Children[i],key,value)}
func(b *BTree)splitChild(parent *node,i int){full:=parent.Children[i];mid:=(b.Order-1)/2;right:=&node{Leaf:full.Leaf};if full.Leaf{right.Keys=append(right.Keys,full.Keys[mid:]...);right.Values=append(right.Values,full.Values[mid:]...);full.Keys=full.Keys[:mid];full.Values=full.Values[:mid];right.Next=full.Next;full.Next=right;parent.Keys=append(parent.Keys,"");copy(parent.Keys[i+1:],parent.Keys[i:]);parent.Keys[i]=right.Keys[0];parent.Children=append(parent.Children,nil);copy(parent.Children[i+2:],parent.Children[i+1:]);parent.Children[i+1]=right;return};promote:=full.Keys[mid];right.Keys=append(right.Keys,full.Keys[mid+1:]...);right.Children=append(right.Children,full.Children[mid+1:]...);full.Keys=full.Keys[:mid];full.Children=full.Children[:mid+1];parent.Keys=append(parent.Keys,"");copy(parent.Keys[i+1:],parent.Keys[i:]);parent.Keys[i]=promote;parent.Children=append(parent.Children,nil);copy(parent.Children[i+2:],parent.Children[i+1:]);parent.Children[i+1]=right}
func(b *BTree)Delete(key string)bool{b.mu.Lock();defer b.mu.Unlock();n:=b.Root;for !n.Leaf{i:=sort.SearchStrings(n.Keys,key);n=n.Children[i]};i:=sort.SearchStrings(n.Keys,key);if i>=len(n.Keys)||n.Keys[i]!=key{return false};n.Keys=append(n.Keys[:i],n.Keys[i+1:]...);n.Values=append(n.Values[:i],n.Values[i+1:]...);return true}
func(b *BTree)Entries()[]Entry{b.mu.RLock();defer b.mu.RUnlock();var out []Entry;var walk func(*node);walk=func(n *node){if n.Leaf{for i,k:=range n.Keys{out=append(out,Entry{k,n.Values[i]})};return};for _,c:=range n.Children{walk(c)}};walk(b.Root);return out}
