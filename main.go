package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/rs/xid"
)

func main() {
	log.SetFlags(log.Llongfile | log.Ldate | log.Lmicroseconds)
	log.Println("Hello")

	//xidMap := make(map[string]bool)

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			http.ServeFile(w, r, "./static/auth.html")
		} else if r.Method == "POST" {
			r.ParseMultipartForm(32 << 20)
			log.Println("logins", r.MultipartForm.Value["user"])
			log.Println("passwords", r.MultipartForm.Value["password"])
			log.Println("map", r.MultipartForm.Value)
			if r.MultipartForm.Value["user"][0] == "qwerty" && r.MultipartForm.Value["password"][0] == "12345" {
				token := xid.New().String()
				//				xidMap[token] = true
				_, err := conn.Do("set", token, "", "EX", 100)
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

	mux.HandleFunc("/suxess", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			cookie, errCookie := r.Cookie("auth")
			if err != nil {
				log.Fatal(err)
			}
			exists, err := redis.Bool(conn.Do("exists", cookie.Value))
			log.Println("Redis get", cookie.Value, exists)
			if err != nil {
				log.Fatal(err)
			}
			//if cookie, err := r.Cookie("auth"); err == nil && xidMap[cookie.Value] {
			if errCookie == nil && exists {
				http.ServeFile(w, r, "./static/suxess.html")
			} else {
				w.Header().Set("Content-Type", "text/html")
				fmt.Fprintln(w, "<h1><pre>(_O_) Go fuck yourself (_*_)")
			}
		} else if r.Method == "POST" {
			//			r.ParseMultipartForm(32 << 20)
			file, handler, err := r.FormFile("uploadfile")
			if err != nil {
				log.Fatal(err)
				return
			}
			defer file.Close()
			f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				panic(err)
				return
			}
			log.Println("Open file", "./upload/"+handler.Filename)
			defer f.Close()
			_, err = io.Copy(f, file)
			if err != nil {
				panic(err)
				return
			}
			http.Redirect(w, r, "http://localhost:9999/suxess", http.StatusFound)
		}
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
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
				MaxAge:   -1,
			})
			_, err = conn.Do("del", cookie.Value)
			log.Println("Redis del", cookie.Value)
			if err != nil {
				log.Fatal(err)
			}
		}
		http.Redirect(w, r, "http://localhost:9999/auth", http.StatusFound)
		//http.ServeFile(w, r, "./static/auth.html")
	})

	mux.Handle("/fileserver/", http.StripPrefix("/fileserver", http.FileServer(http.Dir("./upload/"))))

	//	h.HandleFunc("/fileserver", func(w http.ResponseWriter, r *http.Request) {
	//http.FileServer(http.Dir("/upload")
	//	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "You're lost, go home")
	})
	err = http.ListenAndServe(":9999", mux)
	log.Println("Error listen and serve", err)
}
