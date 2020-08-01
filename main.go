package main

import (
	"fmt"
	"net/http"

	"github.com/odmishien/go-auth-playground/auth"
)

func main() {
	http.Handle("/public", publicHandler)
	http.Handle("/private", auth.JwtMiddleware.Handler(privateHandler))
	http.Handle("/auth", auth.GetTokenHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic("can't ListenAndServe")
	}
}

var publicHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is Public endpoint.")
})

var privateHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is Private endpoint")
})
