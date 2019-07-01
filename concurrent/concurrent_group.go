package concurrent

import (
	"errors"
	"runtime"
	"sort"
	"sync"
)

type STATUS int

const START = 1

type ConcurrentInterface interface {
	input(...interface{}) interface{}
	output(interface{})
}

func NewCPUConcurrentMonitor() *ConcurrentMonitor {
	workerNum := runtime.NumCPU()
	cm, _ := NewConcurrentMonitor(workerNum)
	return cm
}

func NewConcurrentMonitor(workerNum int) (*ConcurrentMonitor, error) {
	if workerNum <= 0 {
		return nil, errors.New("Worker num cannot under 0! ")
	}
	return &ConcurrentMonitor{pipeCtrl: make(chan func(), workerNum), workerNum: workerNum}, nil
}

type Result struct {
	order  int
	result interface{}
}

type ResultSet []*Result

func (rs ResultSet) Len() int {
	return len(rs)
}

func (rs ResultSet) Less(i, j int) bool {
	return rs[i].order < rs[j].order
}

func (rs ResultSet) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

type ConcurrentMonitor struct {
	sync.WaitGroup
	lock      sync.Mutex
	pipeCtrl  chan func()
	status    STATUS
	workerNum int
	resultSet ResultSet
	index     int
}

func (cm *ConcurrentMonitor) GetResultSet() []interface{} {
	sort.Stable(cm.resultSet)
	var result []interface{}
	for _, r := range cm.resultSet {
		result = append(result, r.result)
	}
	return result
}

func (cm *ConcurrentMonitor) init() {
	if cm.status == START {
		return
	}
	defer func() {
		cm.status = START
	}()
	for i := 0; i < cm.workerNum; i++ {
		go func() {
			for callback := range cm.pipeCtrl {
				callback()
			}
		}()
	}
}

// 饿汉并发器，生产消费者模式。适合大规模并发，用于实现了接口对象并发
func (cm *ConcurrentMonitor) HungryConcurrent(process ConcurrentInterface, args ...interface{}) {
	cm.init()
	cm.Add(1)
	cm.pipeCtrl <- func() {
		defer cm.Done()
		process.output(process.input(args...))
	}
}

// 懒汉并发器，每次都会新建协程，适合小规模并发，用于实现了接口对象并发
func (cm *ConcurrentMonitor) LazyConcurrent(process ConcurrentInterface, args ...interface{}) {
	cm.pipeCtrl <- nil
	cm.Add(1)
	go func() {
		defer cm.Done()
		defer func() {
			<-cm.pipeCtrl
		}()
		process.output(process.input(args...))
	}()
}

// 饿汉并发器，生产消费者模式。适合大规模并发，用于函数并发
func (cm *ConcurrentMonitor) HungryConcurrentDirect(fun func() interface{}) {
	cm.init()
	cm.Add(1)
	cm.pipeCtrl <- func(i int) func() {
		return func() {
			defer cm.Done()
			cm.lock.Lock()
			cm.resultSet = append(cm.resultSet, &Result{order: i, result: fun()})
			cm.lock.Unlock()
		}
	}(cm.index)
	cm.index++
}

// 懒汉并发器，每次都会新建协程，适合小规模并发，用于函数并发
func (cm *ConcurrentMonitor) LazyConcurrentDirect(fun func() interface{}) {
	cm.pipeCtrl <- nil
	cm.Add(1)
	go func(i int) func() {
		return func() {
			defer cm.Done()
			defer func() {
				<-cm.pipeCtrl
			}()
			cm.lock.Lock()
			cm.resultSet = append(cm.resultSet, &Result{order: i, result: fun()})
			cm.lock.Unlock()
		}
	}(cm.index)()
	cm.index++
}

func (cm *ConcurrentMonitor) Wait() {
	defer close(cm.pipeCtrl)
	cm.WaitGroup.Wait()
}
