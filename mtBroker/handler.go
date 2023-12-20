package mtbroker

import (
	"math/rand"
	"time"
)

type Handler struct {
	id int
	f  *func(string) string
}

func NewHandler(id int, f func(string) string) *Handler {
	return &Handler{
		id: id,
		f:  &f,
	}
}

func (h *Handler) handle(s string, c chan<- *Handler, r chan<- string) {
	defer func() {
		c <- h
	}()

	time.Sleep(time.Duration(100+rand.Int()%301) * time.Millisecond)
	r <- (*h.f)(s)
}
