package concurrent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	result int
}

func (a *A) input(args ...interface{}) interface{} {
	b, _ := args[0].(int)
	return a.fun(b)
}

func (a *A) output(result interface{}) {
	a.result, _ = result.(int)
}

func (a *A) fun(b int) int {
	return b * 4
}

func fun(b int) int {
	return b * 4
}

func TestConcurrentMonitor_HungryConcurrent(t *testing.T) {
	cm := NewCPUConcurrentMonitor()
	var aList []*A

	for i := 0; i <= 10; i++ {
		aList = append(aList, &A{})
	}

	for i := 0; i <= 10; i++ {
		cm.HungryConcurrent(aList[i], i*3)
	}

	cm.Wait()

	for i := 0; i <= 10; i++ {
		assert.Equal(t, i*12, aList[i].result)
	}
}

func TestConcurrentMonitor_LazyConcurrent(t *testing.T) {
	cm := NewCPUConcurrentMonitor()
	var aList []*A

	for i := 0; i <= 10; i++ {
		aList = append(aList, &A{})
	}

	for i := 0; i <= 10; i++ {
		cm.LazyConcurrent(aList[i], i*3)
	}

	cm.Wait()

	for i := 0; i <= 10; i++ {
		assert.Equal(t, i*12, aList[i].result)
	}
}

func TestConcurrentMonitor_HungryConcurrentDirect(t *testing.T) {
	cm := NewCPUConcurrentMonitor()

	for i := 0; i <= 10; i++ {
		cm.HungryConcurrentDirect(func(b int) func() interface{} {
			return func() interface{} { return fun(b) }
		}(i * 3))
	}

	cm.Wait()

	rs := cm.GetResultSet()
	for i := 0; i <= 10; i++ {
		assert.Equal(t, i*12, rs[i])
	}
}

func TestConcurrentMonitor_LazyConcurrentDirect(t *testing.T) {
	cm := NewCPUConcurrentMonitor()

	for i := 0; i <= 10; i++ {
		cm.LazyConcurrentDirect(func(b int) func() interface{} {
			return func() interface{} { return fun(b) }
		}(i * 3))
	}

	cm.Wait()

	rs := cm.GetResultSet()
	for i := 0; i <= 10; i++ {
		assert.Equal(t, i*12, rs[i])
	}
}
