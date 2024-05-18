package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	w.Header().Set("content-type", "text/plain")
	_, _ = w.Write([]byte("Hello World!"))
}
