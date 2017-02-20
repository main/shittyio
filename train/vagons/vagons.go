package vagons

import (
	"fmt"
	"log"
	"net/http"
)

func HoldPanic(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var a = 0
	log.Println(1 / a)

	log.Println("vagon start")
	defer func() {
		log.Println("defer start")
		if recover() != nil {
			log.Println("recover !=  nil")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Panic")
		}
	}()
	next(w, r)
	log.Println("next handler finish")
}
