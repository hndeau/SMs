package main

import (
	"encoding/json"
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

	http.HandleFunc("/forward", forwardRequest)

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/callback.html")) // Assuming the HTML file is named "index.html"
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(currentDir, "html/chat.html")) // Assuming the HTML file is named "index.html"
	})

	// Add the /cognito handler
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		cognitoURL := "https://smessauth.auth.us-east-1.amazoncognito.com/signup?client_id=34mgjfocrlfp3c4ij35qoe8d4b&response_type=token&scope=email+openid+phone&redirect_uri=http%3A%2F%2Flocalhost%2Fcallback"
		http.Redirect(w, r, cognitoURL, http.StatusFound)
	})

	http.ListenAndServe("80", nil)

	log.Fatal(http.ListenAndServe(":80", nil))
}
func forwardRequest(w http.ResponseWriter, r *http.Request) {
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

	var response = ""
	if r.Method == http.MethodPost {
		var payload outgoing
		err = json.Unmarshal(body, &payload)
		if err != nil {
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		_, err := sendMessage(payload, idToken)
		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		conversationID := r.Header.Get("conversation_id")
		timestamp := r.Header.Get("timestamp")
		_, err := getMessages(conversationID, timestamp, idToken)
		if err != nil {
			http.Error(w, "Failed to send message", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
func sendMessage(payload outgoing, idToken string) ([]byte, error) {
	apiEndpoint := "https://58z24w81cl.execute-api.us-east-1.amazonaws.com/test/sendMessage/"

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

func getMessages(conversationID string, timestamp string, idToken string) ([]byte, error) {
	apiEndpoint := "https://58z24w81cl.execute-api.us-east-1.amazonaws.com/test/getMessages/"

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id_token", idToken)
	req.Header.Set("conversation_id", conversationID)
	req.Header.Set("timestamp", timestamp)

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
