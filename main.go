package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
	log.Println("Hello")

	h := http.NewServeMux()

	h.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit foo!")
		r.Body
		io.ReadCloser
	})

	h.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit bar!")
	})

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "You're lost, go home")
	})

	err := http.ListenAndServe(":9999", h)
	log.Fatal(err)
}
