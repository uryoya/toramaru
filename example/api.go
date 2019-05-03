package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

const (
	port string = ":8071"
)

func main() {
	http.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL)
		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(bodyBytes)
	})

	log.Printf("serve on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
