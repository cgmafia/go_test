package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
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

// Generates Token
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	key := hex.EncodeToString(b)
	//fmt.Println(url.PathEscape(key))
	return key
}

// This is to loop the timer non-stop
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

// Default error 404 page
func Timeoutpage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error 404: Timeout. Key expired")
}

func main() {

	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(30*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handler := http.HandlerFunc(hello)

	go func() {
		if is := Timerloop(); is == false {
			// generate a unique token path
			newurl := "/" + GenerateSecureToken(5)
			fmt.Println(newurl)
			go Timerloop()
		}
	}()

	// http.HandleFunc("/", hello)
	http.Handle("/", cacheClient.Middleware(handler))
	http.HandleFunc("/404", Timeoutpage)
	fmt.Printf("Starting server for testing HTTP POST... \n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
