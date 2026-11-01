// Command metrics-transformer stands in for
// quay.io/bentoml/yatai-bento-metrics-transformer, which has no public source
// and no arm64 image. yatai-deployment runs it next to every Bento container
// (controllers/bentodeployment_controller.go in v1.1.21) with:
//
//	BENTOML_SERVER_HOST, BENTOML_SERVER_PORT  where the Bento serves /metrics
//	PORT                                      where this server listens
//	OLD_METRICS_PREFIX, NEW_METRICS_PREFIX     the metric name prefix to rename
//
// It answers the probes on /healthz and serves the Bento's /metrics with the
// metric names that start with OLD_METRICS_PREFIX renamed to NEW_METRICS_PREFIX.
package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// rename replaces the prefix of metric names in a Prometheus text exposition:
// sample lines that start with the prefix, and the name in # HELP and # TYPE.
func rename(body, oldPrefix, newPrefix string) string {
	if oldPrefix == "" || oldPrefix == newPrefix {
		return body
	}
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		for _, comment := range []string{"# HELP ", "# TYPE "} {
			if strings.HasPrefix(line, comment+oldPrefix) {
				line = comment + newPrefix + line[len(comment)+len(oldPrefix):]
			}
		}
		if strings.HasPrefix(line, oldPrefix) {
			line = newPrefix + line[len(oldPrefix):]
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func metricsHandler(upstream, oldPrefix, newPrefix string) http.HandlerFunc {
	client := &http.Client{Timeout: 10 * time.Second}
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.Get(upstream)
		if err != nil {
			http.Error(w, "bento metrics: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "bento metrics: "+err.Error(), http.StatusBadGateway)
			return
		}
		if ct := resp.Header.Get("Content-Type"); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
		w.WriteHeader(resp.StatusCode)
		io.WriteString(w, rename(string(body), oldPrefix, newPrefix))
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	upstream := "http://" + net.JoinHostPort(getenv("BENTOML_SERVER_HOST", "localhost"), getenv("BENTOML_SERVER_PORT", "3000")) + "/metrics"
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok") })
	mux.Handle("/metrics", metricsHandler(upstream, os.Getenv("OLD_METRICS_PREFIX"), os.Getenv("NEW_METRICS_PREFIX")))
	addr := ":" + getenv("PORT", "3001")
	log.Printf("serving %s as http://%s/metrics", upstream, addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
