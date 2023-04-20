package main

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"
)

var appClientID = ""

func main() {
	// Get the current file's path
	_, currentFilePath, _, _ := runtime.Caller(0)

	// Get the directory of the current file
	currentDir := filepath.Dir(currentFilePath)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the "code" URL parameter
		code := r.URL.Query().Get("code")

		if code != "" {
			// If "code" is provided, create a cookie named "token" with the value of "code"
			token := &http.Cookie{
				Name:    "token",
				Value:   code,
				Expires: time.Now().Add(30 * 24 * time.Hour), // Cookie expires in 30 days
				Path:    "/",
			}
			http.SetCookie(w, token)
		}

		http.ServeFile(w, r, filepath.Join(currentDir, "landingpage.html")) // Assuming the HTML file is named "index.html"
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "chat.html")) // Assuming the HTML file is named "index.html"
	})

	http.ListenAndServe(":8080", nil)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
