package handler

import "net/http"

type Health struct{}

func (h *Health) Handle(
	rw http.ResponseWriter,
	req *http.Request,
) {
	rw.Write([]byte("ok"))
}
