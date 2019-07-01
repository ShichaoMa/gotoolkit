package concurrent

import (
	"testing"
)

func BenchmarkConcurrentMonitor_LazyConcurrent(b *testing.B) {
	cm := NewCPUConcurrentMonitor()
	var aList []*A

	for i := 0; i <= b.N; i++ {
		aList = append(aList, &A{})
	}

	for i := 0; i <= b.N; i++ {
		cm.LazyConcurrent(aList[i], i*3)
	}

	cm.Wait()
}

func BenchmarkConcurrentMonitor_HungryConcurrent(b *testing.B) {
	cm := NewCPUConcurrentMonitor()
	var aList []*A

	for i := 0; i <= b.N; i++ {
		aList = append(aList, &A{})
	}

	for i := 0; i <= b.N; i++ {
		cm.HungryConcurrent(aList[i], i*3)
	}

	cm.Wait()
}

func BenchmarkConcurrentMonitor_HungryConcurrentDirect(b *testing.B) {
	cm := NewCPUConcurrentMonitor()

	for i := 0; i <= b.N; i++ {
		cm.HungryConcurrentDirect(func(b int) func() interface{} {
			return func() interface{} { return fun(b) }
		}(i * 3))
	}

	cm.Wait()
}

func BenchmarkConcurrentMonitor_LazyConcurrentDirect(b *testing.B) {
	cm := NewCPUConcurrentMonitor()

	for i := 0; i <= b.N; i++ {
		cm.LazyConcurrentDirect(func(b int) func() interface{} {
			return func() interface{} { return fun(b) }
		}(i * 3))
	}

	cm.Wait()
}
