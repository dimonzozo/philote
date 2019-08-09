package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

func TestHiveSuccessfulPhiloteRegistration(t *testing.T) {
	h := NewHive()
	if h.PhilotesCount() != 0 {
		t.Error("new Hive shouldn't have registered philotes")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"read":  []string{"test-channel"},
		"write": []string{"test-channel"},
	})

	tokenString, err := token.SignedString(Config.jwtSecret)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
	header := map[string][]string{
		"Authorization": {"Bearer " + tokenString},
	}
	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	_, _, err = websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		t.Error(err)
	}

	// wait for connection message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 1 {
		t.Error("philote should  be registered on successful auth")
	}
}

func TestHiveSuccessfulPhiloteRegistrationWithQuerystring(t *testing.T) {
	h := NewHive()
	if h.PhilotesCount() != 0 {
		t.Error("new Hive shouldn't have registered philotes")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"read":  []string{"test-channel"},
		"write": []string{"test-channel"},
	})

	tokenString, err := token.SignedString(Config.jwtSecret)
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))

	u, _ := url.Parse(server.URL)
	q := u.Query()
	q.Set("auth", tokenString)
	u.RawQuery = q.Encode()
	u.Scheme = "ws"

	_, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		t.Error(err)
	}

	// wait for connection message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 1 {
		t.Error("philote should  be registered on successful auth")
	}
}

func TestHiveIncorrectAuth(t *testing.T) {
	h := NewHive()
	if h.PhilotesCount() != 0 {
		t.Error("new Hive shouldn't have registered philotes")
	}

	server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
	header := map[string][]string{
		"Authorization": {"Bearer " + "foo"},
	}
	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	_, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err == nil {
		t.Error("The Dial action should fail when there is no auth token")
	}

	// wait for connection message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 0 {
		t.Error("philote should not be registered when missing auth")
	}
}

func TestHivePhiloteRegistrationWithNoAuth(t *testing.T) {
	h := NewHive()
	if h.PhilotesCount() != 0 {
		t.Error("new Hive shouldn't have registered philotes")
	}

	server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err == nil {
		t.Error("The Dial action should fail when there is no auth token")
	}

	// wait for connection message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 0 {
		t.Error("philote should not be registered when missing auth")
	}
}

func TestHiveDeregisterPhilote(t *testing.T) {
	h := NewHive()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"read":  []string{"test-channel"},
		"write": []string{"test-channel"},
	})
	tokenString, err := token.SignedString(Config.jwtSecret)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(h.ServeNewConnection))
	header := map[string][]string{
		"Authorization": {"Bearer " + tokenString},
	}
	u, _ := url.Parse(server.URL)
	u.Scheme = "ws"
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		t.Error(err)
	}

	// wait for connection message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 1 {
		t.Error("philote should  be registered on successful auth")
	}

	conn.Close()
	// wait for deregister message to be processed
	time.Sleep(time.Millisecond * 500)

	if h.PhilotesCount() != 0 {
		t.Error("Disconnected Philotes should be automatically deregistered")
	}
}
