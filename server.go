package main

import (
	"fmt"
	"github.com/learn_to_build_web_application_with_go/session"
	_ "github.com/learn_to_build_web_application_with_go/session/provider"
	"github.com/learn_to_build_web_application_with_go/session/sessionMain"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>hello, go!</h1>")
}

func main() {

	register()
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func register() {
	http.HandleFunc("/", index)
	http.HandleFunc("/set_cookie", session.CookieSet)
	http.HandleFunc("/get_cookie", session.CookieGet)
	http.HandleFunc("/go_session", sessionMain.Count)
}
