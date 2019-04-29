package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
)

func main() {
	director := func(request *http.Request) {
		log.Println(request.URL.Path)
		request.URL.Scheme = "http"
		matched, err := regexp.MatchString(`/api/.*`, request.URL.Path)
		if err != nil {
			panic(err)
		}
		if matched {
			request.URL.Host = "localhost:8071"
		} else {
			request.URL.Host = "localhost:8070"
		}
	}
	rp := &httputil.ReverseProxy{Director: director}
	server := http.Server{
		Addr:    ":8080",
		Handler: rp,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
