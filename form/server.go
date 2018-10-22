package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, err := template.ParseFiles("form/login_form.gtpl")
		if err == nil {
			t.Execute(w, nil)
		} else {
			log.Println("login get", err)
			return
		}
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		username := r.Form["username"]
		password := r.FormValue("password")

		switch {
		case len(username) == 0 || username[0] == "":
			fmt.Fprintln(w, "Username is null! Not Allow!")
		case len(password) == 0:
			fmt.Fprintln(w, "Password is null , Not Allow!")
		case username[0] == "zoujiejun" && password == "123456":
			fmt.Fprintln(w, "login success")
		case username[0] != "zoujiejun" || password != "123456":
			fmt.Fprintln(w, "username or password is wrong")
		default:
			fmt.Fprintln(w, "unknown error!")
		}
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//Create a token
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		tmp, err := template.ParseFiles("form/upload.gtpl")
		if err != nil {
			log.Println("upload error:", err)
			return
		} else {
			tmp.Execute(w, token)
		}

	case http.MethodPost:

		/*
			处理文件上传需要调用r.ParseMultipartForm，
			里面的参数表示maxMemory，调用ParseMultipartForm之后，
			上传的文件存储在maxMemory大小的内存里面，
			如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中。
			可以通过r.FormFile获取上面的文件句柄，然后实例中使用了io.Copy来存储文件。
		*/
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println("upload error:", err)
			return
		}
		defer file.Close()
		fmt.Fprintln(w, handler.Header)
		f, err := os.OpenFile("form/MyFile/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("upload error:", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		/*
			上传文件主要三步处理：
			1.表单中增加enctype="multipart/form-data"
			2.服务端调用r.ParseMultipartForm,把上传的文件存储在内存和临时文件中
			3.使用r.FormFile获取文件句柄，然后对文件进行存储等处理。
		*/
	default:
		fmt.Fprintln(w, "404 not found!")

	}
}

func main() {
	http.HandleFunc("/login", login)

	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":9512", nil)
	if err != nil {
		log.Fatal(err)
	}
}
