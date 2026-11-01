package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const exposition = `# HELP BENTOML_capstone:url_request_total Total requests
# TYPE BENTOML_capstone:url_request_total counter
BENTOML_capstone:url_request_total{endpoint="/api/predict"} 3
bentoml_api_server_request_in_progress{endpoint="/api/predict"} 0
`

func TestRename(t *testing.T) {
	got := rename(exposition, "BENTOML_capstone:url_", "BENTOML_")
	want := `# HELP BENTOML_request_total Total requests
# TYPE BENTOML_request_total counter
BENTOML_request_total{endpoint="/api/predict"} 3
bentoml_api_server_request_in_progress{endpoint="/api/predict"} 0
`
	if got != want {
		t.Fatalf("rename:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenameWithoutPrefixKeepsBody(t *testing.T) {
	if got := rename(exposition, "", "BENTOML_"); got != exposition {
		t.Fatalf("an empty old prefix changed the body:\n%s", got)
	}
}

func TestMetricsHandler(t *testing.T) {
	bento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		io.WriteString(w, exposition)
	}))
	defer bento.Close()

	rec := httptest.NewRecorder()
	metricsHandler(bento.URL+"/metrics", "BENTOML_capstone:url_", "BENTOML_")(rec, httptest.NewRequest("GET", "/metrics", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; version=0.0.4" {
		t.Fatalf("content type %q", ct)
	}
	if got := rec.Body.String(); got != rename(exposition, "BENTOML_capstone:url_", "BENTOML_") {
		t.Fatalf("body:\n%s", got)
	}
}

func TestMetricsHandlerBentoDown(t *testing.T) {
	rec := httptest.NewRecorder()
	metricsHandler("http://127.0.0.1:1/metrics", "", "")(rec, httptest.NewRequest("GET", "/metrics", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status %d, want 502", rec.Code)
	}
}
