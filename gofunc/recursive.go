package gofunc

import (
	"context"
	"sync"
	"sync/atomic"
)

type Recursive struct {
	doCount int

	chInput  chan interface{}
	chBuffer chan interface{}
	buffer   []interface{}

	numOfJob int64
	doneJob  chan struct{}

	Do       func(input interface{}, output chan<- interface{}, recursiveFunc func(v interface{}), v ...interface{})
	ChOutput chan<- interface{}

	metricRunningCount int64
}

func New() Recursive {
	return Recursive{doCount: 100, chInput: make(chan interface{}), chBuffer: make(chan interface{}),
		buffer: make([]interface{}, 0), doneJob: make(chan struct{})}
}

func (c *Recursive) AddInput(input string) {
	atomic.AddInt64(&c.numOfJob, 1)
	c.buffer = append(c.buffer, input)
}

func (c *Recursive) Run(v ...interface{}) {
	var wgPool, wgCtl sync.WaitGroup
	ctx, ctxCancel := context.WithCancel(context.Background())

	wgCtl.Add(1)
	go func() {
		defer wgCtl.Done()

		for ctx.Err() == nil {
			if len(c.buffer) == 0 {
				select {
				case in := <-c.chInput:
					c.buffer = append(c.buffer, in)
				case <-ctx.Done():
					break
				}
			} else {
				select {
				case in := <-c.chInput:
					c.buffer = append(c.buffer, in)
				case c.chBuffer <- c.buffer[0]:
					c.buffer = c.buffer[1:]
				case <-ctx.Done():
					break
				}
			}
		}
	}()

	wgPool.Add(c.doCount)
	for i := 0; i < c.doCount; i++ {
		go func() {
			defer wgPool.Done()
			for value := range c.chBuffer {
				atomic.AddInt64(&c.metricRunningCount, 1)
				c.Do(value, c.ChOutput, c.RecursiveInput, v...)
				c.doneJob <- struct{}{}
			}
		}()
	}

	go func() {
		defer close(c.ChOutput)
		defer wgPool.Wait()
		for ctx.Err() == nil {
			<-c.doneJob
			if atomic.AddInt64(&c.numOfJob, -1) == 0 {
				ctxCancel()
				wgCtl.Wait()
				close(c.chInput)
				close(c.chBuffer)
			}
		}
	}()
}

func (c *Recursive) RecursiveInput(v interface{}) {
	atomic.AddInt64(&c.numOfJob, 1)
	c.chInput <- v
}

func (c *Recursive) GetRunningCount() int64 {
	return atomic.LoadInt64(&c.metricRunningCount)
}
