package mtcounter

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"

	"time"
)

type webCounter struct {
	visits *sync.Map
}

func StartConcurrentVisits(done chan<- struct{}) {
	wc := newWebCounter()

	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go visitor(i, 100, wc, &wg, "google.com", "amazon.com", "go.dev", "apple.com")
	}

	wg.Wait()

	for _, s := range []string{"google.com", "amazon.com", "go.dev", "apple.com", "spotify.com"} {
		v, err := wc.GetVisitors(s)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("%s visited %d times\n", s, v)
	}

	done <- struct{}{}
}

func newWebCounter() *webCounter {
	return &webCounter{visits: &sync.Map{}}
}

func (wc *webCounter) Increment(url string) {
	val, ok := wc.visits.LoadOrStore(url, int32(0))
	if ok {
		current := val.(int32)
		atomic.AddInt32(&current, 1)
		wc.visits.Store(url, current)
	}
}

func (wc *webCounter) GetVisitors(url string) (int, error) {
	val, ok := wc.visits.Load(url)
	if ok {
		current := val.(int32)
		return int(atomic.LoadInt32(&current)), nil
	}

	return -1, fmt.Errorf("failed to load visits for : %s", url)
}

func visitor(id int, visits int, wc *webCounter, wg *sync.WaitGroup, urls ...string) {
	for i := 0; i < visits; i++ {
		timeout := time.Duration(100 + (rand.Int() % 401))
		time.Sleep(timeout * time.Millisecond)

		url := urls[rand.Int()%len(urls)]
		wc.Increment(url)

		fmt.Printf("visitor-%d visited %s\n", id, url)
	}

	wg.Done()
}
