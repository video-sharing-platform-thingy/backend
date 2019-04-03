package main

import (
	"fmt"
	"log"
	"net/http"
)

// logger is a middleware that simple logs requests.
func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)

		if loadedConfig.LogSession {
			session, err := sessionStore.Get(r, loadedConfig.SessionCookieName)
			if err == nil {
				buffer := "Session: "
				first := true
				for k, v := range session.Values {
					if first {
						first = false
					} else {
						buffer += ", "
					}
					buffer += fmt.Sprintf("%s = %v", k, v)
				}
				log.Print(buffer)
			}
		}

		h.ServeHTTP(w, r)
	})
}
