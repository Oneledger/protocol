package main

import (
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/Oneledger/protocol/log"
)

type Controller struct {
	functionInit  []func(LoadTestNode, int, *SharedData) (interface{}, error)
	functionExec  []func(LoadTestNode, int, *SharedData, interface{})
	functionindex int
	stopChan      chan bool
	waiter        sync.WaitGroup
	counterChan   chan int
	osChan        chan os.Signal
	logger        *log.Logger
	waitduration  time.Duration
	sharedData    *SharedData
	threadData    []interface{}
}

func NewController(threads int, waitduration time.Duration) *Controller {
	stopChan := make(chan bool, threads)
	waiter := sync.WaitGroup{}
	counterChan := make(chan int, threads)
	c := make(chan os.Signal, 1)
	shared := NewSharedData()

	return &Controller{
		functionInit:  nil,
		functionExec:  nil,
		functionindex: 0,
		stopChan:      stopChan,
		waiter:        waiter,
		counterChan:   counterChan,
		osChan:        c,
		logger:        log.NewLoggerWithPrefix(os.Stdout, "controller"),
		waitduration:  waitduration,
		sharedData:    shared,
		threadData:    nil,
	}
}

func (c *Controller) AddExecutorFunction(fExec func(LoadTestNode, int, *SharedData, interface{}), fInit func(LoadTestNode, int, *SharedData) (interface{}, error)) {
	c.functionExec = append(c.functionExec, fExec)
	c.functionInit = append(c.functionInit, fInit)
}

func (c *Controller) InitThread(node LoadTestNode, thread int) {
	findex := c.functionindex
	fInit := c.functionInit[findex]
	if fInit != nil {
		out, err := fInit(node, thread, c.sharedData)
		if err != nil {
			logger.Fatalf("Error initializing thread data | Thread no :%d | %s", thread, err.Error())
		}
		c.threadData = append(c.threadData, out)
	} else {
		c.threadData = append(c.threadData, nil)
	}

	c.functionindex++
	if c.functionindex == len(c.functionExec) {
		c.functionindex = 0
	}
}

func (c *Controller) AddThread(node LoadTestNode, thread int) {
	if len(c.functionExec) == 0 {
		logger.Fatal("Function List is empty")
	}
	c.waiter.Add(1)
	findex := c.functionindex
	go func(stop chan bool) {
		name, _ := node.clCtx.FullNodeClient().NodeName()
		logger.Infof("Starting Broadcast for: %s | Thread no :%d | Node-Name :%s | Tx-Type : %s", node.superadminName, thread, name.Name, GetFunctionName(c.functionExec[findex]))
		for {
			in := c.threadData[thread]
			c.functionExec[findex](node, thread, c.sharedData, in)
			c.counterChan <- 1
			select {
			//stop populated true by iterating over al the threads
			case <-stop:
				c.waiter.Done()
				return
			default:
				time.Sleep(c.waitduration)
			}
		}
	}(c.stopChan)
	c.functionindex++
	if c.functionindex == len(c.functionExec) {
		c.functionindex = 0
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
