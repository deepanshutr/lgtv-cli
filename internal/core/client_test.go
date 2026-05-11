package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_VolumeAbsolute(t *testing.T) {
	var got map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/volume" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(b, &got)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	if err := New(srv.URL).VolumeAbsolute(context.Background(), 25); err != nil {
		t.Fatal(err)
	}
	if got["level"].(float64) != 25 {
		t.Errorf("expected level=25, got %v", got)
	}
}

func TestClient_State(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"ok":true,"volume":12,"apps":["netflix"]}`))
	}))
	defer srv.Close()

	got, err := New(srv.URL).State(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got["volume"].(float64) != 12 {
		t.Errorf("expected volume=12, got %v", got)
	}
}

func TestClient_ErrorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(502)
		w.Write([]byte(`{"detail":"TV unreachable"}`))
	}))
	defer srv.Close()

	err := New(srv.URL).PowerOff(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
