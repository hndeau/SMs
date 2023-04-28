package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
)

func serveStaticFile(currentDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(currentDir, r.URL.Path[1:])
		http.ServeFile(w, r, path)
	})
}

type outgoing struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
}

type incoming struct {
	ConversationID string `json:"conversation_id"`
	Timestamp      int    `json:"timestamp"`
}

func main() {
	// Get the current file's path
	_, currentFilePath, _, _ := runtime.Caller(0)

	// Get the directory of the current file
	currentDir := filepath.Dir(currentFilePath)

	// Serve static files
	http.Handle("/js/", serveStaticFile(currentDir))
	http.Handle("/css/", serveStaticFile(currentDir))
	http.Handle("/html/", serveStaticFile(currentDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/landingpage.html")) // Assuming the HTML file is named "index.html"
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/login.html")) // Assuming the HTML file is named "index.html"
	})

	http.HandleFunc("/send", sendMessageHandler)

	http.HandleFunc("/retrieve", getMessageHandler)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/callback.html")) // Assuming the HTML file is named "index.html"
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/chat.html")) // Assuming the HTML file is named "index.html"
	})

	// Add the /cognito handler
	//http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	//	cognitoURL := ""
	//	http.Redirect(w, r, cognitoURL, http.StatusFound)
	//})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/signup.html")) // Assuming the HTML file is named "index.html"
	})

	http.ListenAndServe("80", nil)

	log.Fatal(http.ListenAndServe(":80", nil))
}
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("id_token")
	if idToken == "" {
		http.Error(w, "Missing id_token", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var payload outgoing
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	response, err := sendMessage(payload, idToken)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("id_token")
	if idToken == "" {
		http.Error(w, "Missing id_token", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var payload incoming
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	response, err := getMessages(payload, idToken)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Print the request header and body to the console
	fmt.Printf("Request Header: %v\n", r.Header)
	fmt.Printf("Request Body: %s\n", string(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
func sendMessage(payload outgoing, idToken string) ([]byte, error) {
	apiEndpoint := "https://58z24w81cl.execute-api.us-east-1.amazonaws.com/prod/sendMessage"

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(string(payloadJSON)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id_token", idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func getMessages(payload incoming, idToken string) ([]byte, error) {
	apiEndpoint := "https://58z24w81cl.execute-api.us-east-1.amazonaws.com/prod/getMessages"
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiEndpoint, strings.NewReader(string(payloadJSON)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id_token", idToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}
