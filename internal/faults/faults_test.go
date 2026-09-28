package faults
import "testing"
func TestInjector(t *testing.T){var i Injector;i.FailAfter(2);if i.Check()!=nil{t.Fatal("failed too early")};if e:=i.Check();e!=ErrInjected{t.Fatal("missing injected error")}}
