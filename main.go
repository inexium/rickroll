package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

// define a page structure
type Page struct {
	Title string
	Body  []byte
}

// TODO : check if this has any usecase now
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// the handler when accessing /
func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	t, _ := template.ParseFiles("./templates/base.html", "./templates/main.html")
	t.Execute(w, p)
}

// clientIP extracts the requester's IP address, stripping the port if present
func clientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// logRequests logs the client IP and access timestamp for every request to stdout
func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s - %s - %s %s\n", clientIP(r), time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// TODO: create a contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	/*title := r.URL.Path[len("/contact/"):]
	p, err := loadPage("./templates/base")
	if err != nil {
		p = &Page{Title: title}
	}*/
	p := &Page{Title: "random"}
	t, _ := template.ParseFiles("./templates/base.html", "./templates/main.html")
	t.Execute(w, p)
}

func main() {
	mux := http.NewServeMux()

	// load static path
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	// define a handler for /
	mux.HandleFunc("/", homeHandler)

	// TODO : check contact handler
	//mux.HandleFunc("/contact", contactHandler)

	// serve, logging the IP and timestamp of every access
	log.Fatal(http.ListenAndServe(":8080", logRequests(mux)))
}
