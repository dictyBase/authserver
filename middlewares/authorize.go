package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func AuthorizeMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		hdr := r.Header
		if hdr.Get("X-Scheme") != "https" {
			http.Error(
				w,
				fmt.Sprintf("scheme is %s not https", hdr.Get("X-Schema")),
				http.StatusBadRequest,
			)
			return
		}
		if hdr.Get("X-Original-Method") == "OPTIONS" {
			w.Write([]byte("passthrough for OPTIONS method"))
			return
		}
		if hdr.Get("X-Original-Method") == "GET" {
			w.Write([]byte("passthrough for GET method"))
			return
		}
		log.Printf("original uri %s\n", hdr.Get("X-Original-Uri"))
		if strings.HasPrefix(hdr.Get("X-Original-Uri"), "/tokens") {
			w.Write([]byte("no validation for /tokens"))
			return
		}
		log.Println("going for JWT check")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
