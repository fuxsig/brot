package wrapper

import "net/http"

type Handler interface {
	ServeChain(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeChain(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type Wrapper struct {
	handler Handler
	next    *Wrapper
}

func (m *Wrapper) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeChain(rw, r, m.next.ServeHTTP)
}

func build(handlers []Handler) *Wrapper {
	if len(handlers) == 0 {
		return &Wrapper{
			HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
			&Wrapper{},
		}
	}
	return &Wrapper{handlers[0], build(handlers[1:])}
}

func New(handlers ...Handler) *Wrapper {
	return build(handlers)
}

func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}
