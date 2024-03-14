package db

import (
	"crypto/sha1"
	"fmt"
	"sync"
	"testing"
)

var cases = []struct {
	runCount int
}{
	{runCount: 100},
	{runCount: 1000},
	{runCount: 10000},
	{runCount: 100000},
}

func task() {
	sha1.Sum([]byte("hello world"))
}
func testChannel(runCount int) {
	doneChan := make(chan struct{}, runCount)

	// Execute tasks
	for i := 0; i < runCount; i++ {
		go func() {
			defer func() { doneChan <- struct{}{} }()
			task()
		}()
	}

	// Wait for all tasks to complete
	for i := 0; i < runCount; i++ {
		<-doneChan
	}
}

func testWaitGroup(runCount int) {
	wg := sync.WaitGroup{}
	wg.Add(runCount)

	// Execute tasks
	for i := 0; i < runCount; i++ {
		go func() {
			defer wg.Done()
			task()
		}()
	}

	// Wait for all tasks to complete
	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	for _, c := range cases {
		b.Run(fmt.Sprintf("runCount%d", c.runCount), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testChannel(c.runCount)
			}
		})
	}
}

func BenchmarkWaitGroup(b *testing.B) {
	for _, c := range cases {
		b.Run(fmt.Sprintf("runCount%d", c.runCount), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testWaitGroup(c.runCount)
			}
		})
	}
}
