package faults

import("errors";"sync")
var ErrInjected=errors.New("injected failure")
type Injector struct{mu sync.Mutex;remaining int}
func(i *Injector)FailAfter(n int){i.mu.Lock();defer i.mu.Unlock();i.remaining=n}
func(i *Injector)Check()error{i.mu.Lock();defer i.mu.Unlock();if i.remaining<=0{return nil};i.remaining--;if i.remaining==0{return ErrInjected};return nil}
