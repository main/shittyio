package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/rs/xid"
)

// TODO: Implement Redis based session storage

func main() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
	log.Println("Hello")

	xidMap := make(map[string]bool)

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	h := http.NewServeMux()

	h.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "./static/auth.html")
		} else if r.Method == "POST" {
			r.ParseMultipartForm(32 << 20)
			log.Println("logins", r.MultipartForm.Value["user"])
			log.Println("passwords", r.MultipartForm.Value["password"])
			log.Println("map", r.MultipartForm.Value)
			if r.MultipartForm.Value["user"][0] == "qwerty" && r.MultipartForm.Value["password"][0] == "12345" {
				token := xid.New().String()
				xidMap[token] = true
				_, err := conn.Do("set", token, "")
				log.Println("Redis set", token)
				if err != nil {
					log.Fatal(err)
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "auth",
					Value:    token,
					Path:     "/",
					Domain:   "localhost",
					HttpOnly: true,
				})
				http.Redirect(w, r, "http://localhost:9999/suxess", http.StatusFound)
			} else {
				http.ServeFile(w, r, "./static/auth.html")
			}
		}
		log.Println("User reached /auth by method", r.Method)
	})

	h.HandleFunc("/suxess", func(w http.ResponseWriter, r *http.Request) {
		cookie, errCookie := r.Cookie("auth")
		if err != nil {
			log.Fatal(err)
		}
		exists, err := conn.Do("exists", cookie.Value)
		exists = redis.Bool()
		log.Println("Redis get", cookie.Value, valRedis, reflect.TypeOf(valRedis))
		if err != nil {
			log.Fatal(err)
		}
		//if cookie, err := r.Cookie("auth"); err == nil && xidMap[cookie.Value] {
		if errCookie == nil && valRedis != nil {
			http.ServeFile(w, r, "./static/suxess.html")
		} else {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, "<h1><pre>(_O_) Go fuck yourself (_*_)")
		}

	})

	h.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		log.Println("User reached /logout by method", r.Method)
		if cookie, err := r.Cookie("auth"); err == nil {
			log.Println("Logout cookie setter")
			http.SetCookie(w, &http.Cookie{
				Name:     "auth",
				Value:    cookie.Value,
				Path:     "/",
				Domain:   "localhost",
				HttpOnly: true,
				Expires:  time.Now(),
			})
			_, err = conn.Do("set", cookie.Value, 0)
			log.Println("Redis set", cookie.Value, 0)
			if err != nil {
				log.Fatal(err)
			}
		}
		http.Redirect(w, r, "http://localhost:9999/auth", http.StatusFound)
		//http.ServeFile(w, r, "./static/auth.html")
	})

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "You're lost, go home")
	})
	err = http.ListenAndServe(":9999", h)
	log.Fatal(err)
}
