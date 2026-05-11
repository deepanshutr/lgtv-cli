package tg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/deepanshutr/lgtv-cli/internal/core"
)

func TestDispatch_UnknownCommand(t *testing.T) {
	r := dispatch(context.Background(), core.New("http://invalid"), "/foo")
	if !strings.Contains(r, "Unknown command") {
		t.Errorf("expected unknown command, got %q", r)
	}
}

func TestDispatch_Vol(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	c := core.New(srv.URL)
	r := dispatch(context.Background(), c, "/vol 25")
	if !strings.Contains(r, "= 25") {
		t.Errorf("expected '= 25', got %q", r)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestDispatch_VolDelta(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()
	r := dispatch(context.Background(), core.New(srv.URL), "/vol +3")
	if !strings.Contains(r, "+3") {
		t.Errorf("expected '+3' in reply, got %q", r)
	}
}
