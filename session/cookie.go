package session

import (
	"fmt"
	"net/http"
	"time"
)

//设置cookie
func CookieSet(w http.ResponseWriter, r *http.Request) {
	//设置到期时间
	expiration := time.Now()
	expiration = expiration.AddDate(0, 0, 7)
	//设置cookie
	cookie := http.Cookie{Name: "age", Value: "21", Expires: expiration}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "name", Value: "zoujiejun", Expires: expiration}
	http.SetCookie(w, &cookie)
	cookie = http.Cookie{Name: "wife", Value: "luxingyu", Expires: expiration}
	http.SetCookie(w, &cookie)
}

//获得cookie
func CookieGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>one by one</h1><br><hr>")
	//单个取值
	name, err := r.Cookie("name")
	if err == nil {
		fmt.Fprintf(w, "<h2>name:%v</h2>", name)
	}

	fmt.Fprintln(w, "<hr><br><h1>range</h1>")
	//循环取值
	for _, cookie := range r.Cookies() {
		fmt.Fprintf(w, "<h2>%v:%v</h2>", cookie.Name, cookie.Value)
	}
}
