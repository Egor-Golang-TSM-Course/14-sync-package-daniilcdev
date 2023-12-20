package mtbroker

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var freeHandlers []*Handler
var mu sync.Mutex

func StartRequestBroking(done chan<- struct{}) {
	for i := 0; i < 10; i++ {
		h := NewHandler(i+1, func(s string) string {
			return strings.Repeat(s, 3)
		})
		freeHandlers = append(freeHandlers, h)
	}

	requests := make(chan string)
	free := make(chan *Handler)

	responses := make(chan string, 100)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer close(responses)
		wg.Wait()
	}()

	go runBroker(requests, free, responses)

	go runEmitter(&wg, requests)

	for s := range responses {
		wg.Done()
		fmt.Printf("response : %s\n", s)
	}

	wg.Wait()
	done <- struct{}{}
}

func runBroker(r chan string, c chan *Handler, resp chan<- string) {
	for {
		select {
		case s, ok := <-r:
			if !ok {
				fmt.Printf("requests channel closed\n")
				return
			}

			mu.Lock()
			h := freeHandlers[0]
			fmt.Printf("handler-%d selected to serve request\n", h.id)
			freeHandlers = freeHandlers[1:]
			go h.handle(s, c, resp)
			mu.Unlock()
		case h := <-c:
			mu.Lock()
			fmt.Printf("handler-%d served request, return to pool\n", h.id)
			freeHandlers = append(freeHandlers, h)
			mu.Unlock()
		default:
			continue
		}
	}
}

func runEmitter(wg *sync.WaitGroup, r chan string) {
	for i := 0; i < 100; i++ {
		wg.Add(1)
		r <- fmt.Sprintf("rs-|%d|", i)
		time.Sleep(50 * time.Millisecond)
	}
	wg.Done()

	close(r)
}
