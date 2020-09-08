package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Println(r.PostForm)
		fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		name := r.FormValue("name")
		message := r.FormValue("message")
		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Message = %s\n", message)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	key := hex.EncodeToString(b)
	//fmt.Println(url.PathEscape(key))
	return key
}

func Timerloop() bool {
	ticker := time.NewTicker(5 * time.Second)
	c := make(chan struct{})
	for {
		select {
		case <-c:
			return false
		case t := <-ticker.C:
			fmt.Println("tick", t)
			//timer is on
			return true
		}
	}

	time.Sleep(3 * time.Millisecond)
	ticker.Stop()
	// fmt.Println("Ticker stopped")
	return false
}

func Timeoutpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error 404: Timeout. Key expired")
}

func main() {
	go func() {
		if is := Timerloop(); is == false {
			// generate a unique token path
			newurl := "/" + GenerateSecureToken(5)
			fmt.Println(newurl)
			go Timerloop()
		}
	}()

	http.HandleFunc("/", hello)
	http.HandleFunc("/404", Timeoutpage)
	fmt.Printf("Starting server for testing HTTP POST... \n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
