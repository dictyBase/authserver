package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

func AuthorizeMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("getting authorized")
		w.Write([]byte("all in good"))
		hdr := r.Header
		if hdr.Get("X-Scheme") != "https" {
			http.Error(
				w,
				fmt.Sprintf("scheme is %s not https", hdr.Get("X-Schema")),
				http.StatusBadRequest,
			)
			return
		}
		if hdr.Get("X-Original-Method") == "GET" {
			w.Write([]byte("passthrough for GET method"))
			return
		}
		if hdr.Get("X-Original-Uri") == "/tokens" {
			w.Write([]byte("no validation for /tokens"))
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
