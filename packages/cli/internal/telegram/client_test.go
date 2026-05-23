package telegram

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("chat_id") != "123" {
			t.Errorf("expected chat_id '123', got '%s'", r.FormValue("chat_id"))
		}
		if r.FormValue("text") != "hello" {
			t.Errorf("expected text 'hello', got '%s'", r.FormValue("text"))
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true,"result":{}}`)
	}))
	defer server.Close()

	client := &Client{HTTPClient: server.Client(), BaseURL: server.URL}
	err := client.Send(SendParams{Token: "test-token", ChatID: "123", Text: "hello"})
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}
}

func TestSend_TelegramAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":false,"description":"Unauthorized"}`)
	}))
	defer server.Close()

	client := &Client{HTTPClient: server.Client(), BaseURL: server.URL}
	err := client.Send(SendParams{Token: "bad-token", ChatID: "123", Text: "hello"})

	if err == nil {
		t.Fatal("expected error for API error response, got nil")
	}
	if err.Error() != "telegram API error: Unauthorized" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestSend_EmptyToken(t *testing.T) {
	client := NewClient()
	err := client.Send(SendParams{Token: "", ChatID: "123", Text: "hello"})

	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestSend_EmptyChatID(t *testing.T) {
	client := NewClient()
	err := client.Send(SendParams{Token: "123:ABC", ChatID: "", Text: "hello"})

	if err == nil {
		t.Fatal("expected error for empty chat ID, got nil")
	}
}

func TestSend_EmptyText(t *testing.T) {
	client := NewClient()
	err := client.Send(SendParams{Token: "123:ABC", ChatID: "123", Text: ""})

	if err == nil {
		t.Fatal("expected error for empty text, got nil")
	}
}

func TestSend_RequestURL(t *testing.T) {
	var receivedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"ok":true,"result":{}}`)
	}))
	defer server.Close()

	client := &Client{HTTPClient: server.Client(), BaseURL: server.URL}
	err := client.Send(SendParams{Token: "123:ABCDEF", ChatID: "42", Text: "test message"})
	if err != nil {
		t.Fatalf("Send() returned error: %v", err)
	}

	expected := "/bot123:ABCDEF/sendMessage"
	if receivedPath != expected {
		t.Errorf("expected path '%s', got '%s'", expected, receivedPath)
	}
}

func TestSend_NetworkError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	server.Close()

	client := &Client{HTTPClient: server.Client(), BaseURL: server.URL}
	err := client.Send(SendParams{Token: "123:ABC", ChatID: "42", Text: "hello"})

	if err == nil {
		t.Fatal("expected error for network failure, got nil")
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.BaseURL != defaultBaseURL {
		t.Errorf("expected BaseURL '%s', got '%s'", defaultBaseURL, client.BaseURL)
	}
	if client.HTTPClient == nil {
		t.Fatal("HTTPClient should not be nil")
	}
}
