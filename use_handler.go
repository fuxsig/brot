package brot

import (
	"net/http"

	"github.com/fuxsig/brot/di"
	"github.com/fuxsig/brot/wrapper"
)

type WrapperHandler struct {
	Wrapper []wrapper.Handler `brot:"wrapper,mandatory"`
	Handler ProvidesHandler   `brot:"handler,mandatory"`
	w       *wrapper.Wrapper
}

func (h *WrapperHandler) InitFunc() (err error) {
	args := make([]wrapper.Handler, len(h.Wrapper)+1)
	copy(args, h.Wrapper)
	args[len(h.Wrapper)] = wrapper.Wrap(h.Handler.HandlerFunc())
	h.w = wrapper.New(args...)
	return
}

func (h *WrapperHandler) Retry() bool {
	return false
}

func (h *WrapperHandler) HandlerFunc() http.Handler {
	return h.w
}

var _ di.ProvidesInit = (*WrapperHandler)(nil)
var _ ProvidesHandler = (*WrapperHandler)(nil)
var _ = di.GlobalScope.Declare((*WrapperHandler)(nil))
