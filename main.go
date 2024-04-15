package main

import (
	"fmt"
	"net/http"
)

var FollowingUser []string

func main() {
	http.HandleFunc("/", helloHandleFunc)
	http.HandleFunc("/follow/:username", followUser)
	http.ListenAndServe(":8080", nil)

}

func helloHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
	fmt.Printf("request made at %s\n", r.URL.Path)
}

func followUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	fmt.Printf("request made at %s\n", r.URL.Path)
	fmt.Fprintf(w, "Hello %s!", username)
}
