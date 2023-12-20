package mtlogger

import (
	"fmt"
	"sync"
)

type LogBuffer struct {
	buffer []string
	mu     sync.Mutex
}

func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		buffer: make([]string, 0, 10),
		mu:     sync.Mutex{},
	}
}

func StartConcurrentLogging(done chan<- struct{}) {
	lb := NewLogBuffer()
	wg := sync.WaitGroup{}
	for i := int64(0); i < 1000; i++ {
		wg.Add(1)

		id := i
		go func(ba *LogBuffer) {
			for k := 0; k < 100; k++ {
				lb.WriteLog(fmt.Sprintf("worker-%d; message %d", id, k))
			}

			wg.Done()
		}(lb)
	}

	wg.Wait()

	for _, msg := range lb.buffer {
		fmt.Printf("msg from buffer - %s\n", msg)
	}

	done <- struct{}{}
}

func (lb *LogBuffer) WriteLog(msg string) {
	defer lb.mu.Unlock()

	lb.mu.Lock()
	lb.buffer = append(lb.buffer, msg)
}
