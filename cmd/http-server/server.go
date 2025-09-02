package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/voidarchive/nepseauth/nepse"
)

type app struct {
	client nepse.Client
}

func main() {
	// Config
	host := getenv("HOST", "127.0.0.1")
	port := getenv("PORT", "8081")
	tlsVerify := getenv("TLS_VERIFY", "false") != "false"

	// Create client
	var (
		c   nepse.Client
		err error
	)
	if tlsVerify {
		c, err = nepse.NewClientWithTLS(true)
	} else {
		c, err = nepse.NewClientWithTLS(false)
	}
	if err != nil {
		log.Fatalf("failed to create nepse client: %v", err)
	}
	defer c.Close(context.Background())

	a := &app{client: c}

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
	})

	// Docs
	mux.HandleFunc("/docs", docsHTML)
	mux.HandleFunc("/docs/openapi.json", docsOpenAPI)

	// Test endpoints (invoke library calls)
	mux.HandleFunc("/test/market/summary", a.handleMarketSummary)
	mux.HandleFunc("/test/market/status", a.handleMarketStatus)
	mux.HandleFunc("/test/top/gainers", a.handleTopGainers)
	mux.HandleFunc("/test/security/", a.handleSecurityRoutes)

	addr := fmt.Sprintf("%s:%s", host, port)
	log.Printf("listening on http://%s (docs: http://%s/docs)", addr, addr)
	if err := http.ListenAndServe(addr, logRequests(mux)); err != nil {
		log.Fatal(err)
	}
}

func (a *app) handleMarketSummary(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	s, err := a.client.GetMarketSummary(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (a *app) handleMarketStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	s, err := a.client.GetMarketStatus(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, s)
}

func (a *app) handleTopGainers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	v, err := a.client.GetTopGainers(ctx)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

// /test/security/{symbol}/...
func (a *app) handleSecurityRoutes(w http.ResponseWriter, r *http.Request) {
	// Expect /test/security/{symbol}/{action}
	// actions: company | depth | history
	p := strings.TrimPrefix(r.URL.Path, "/test/security/")
	parts := strings.Split(p, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	symbol := parts[0]
	action := ""
	if len(parts) > 1 {
		action = parts[1]
	}

	switch action {
	case "company":
		a.handleCompanyBySymbol(w, r, symbol)
	case "depth":
		a.handleDepthBySymbol(w, r, symbol)
	case "history":
		a.handleHistoryBySymbol(w, r, symbol)
	default:
		http.NotFound(w, r)
	}
}

func (a *app) handleCompanyBySymbol(w http.ResponseWriter, r *http.Request, symbol string) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	d, err := a.client.GetCompanyDetailsBySymbol(ctx, symbol)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (a *app) handleDepthBySymbol(w http.ResponseWriter, r *http.Request, symbol string) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	d, err := a.client.GetMarketDepthBySymbol(ctx, symbol)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, d)
}

func (a *app) handleHistoryBySymbol(w http.ResponseWriter, r *http.Request, symbol string) {
	q := r.URL.Query()
	start := q.Get("start")
	end := q.Get("end")
	if start == "" || end == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "start and end are required (YYYY-MM-DD)"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	h, err := a.client.GetPriceVolumeHistoryBySymbol(ctx, symbol, start, end)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, h)
}

// --- Helpers ---

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, err error) {
	type errResp struct {
		Error string `json:"error"`
	}
	writeJSON(w, http.StatusBadGateway, errResp{Error: err.Error()})
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// --- Docs ---

func docsHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io := `<!doctype html>
<html>
<head>
  <meta charset="utf-8"/>
  <title>NEPSE Test Server Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>body{margin:0;}#swagger-ui{max-width:100%;}</style>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '/docs/openapi.json',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: 'BaseLayout'
    });
  }
  </script>
  </head>
  <body>
    <div id="swagger-ui"></div>
  </body>
  </html>`
	_, _ = w.Write([]byte(io))
}

func docsOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(openapiJSON))
}

const openapiJSON = `{
  "openapi": "3.0.3",
  "info": {"title": "NEPSE Test Server (Library-backed)", "version": "1.0.0"},
  "servers": [{"url": "/"}],
  "paths": {
    "/health": {
      "get": {"summary": "Health", "responses": {"200": {"description": "ok"}}}
    },
    "/test/market/summary": {
      "get": {"summary": "Market summary", "responses": {"200": {"description": "OK"}}}
    },
    "/test/market/status": {
      "get": {"summary": "Market status", "responses": {"200": {"description": "OK"}}}
    },
    "/test/top/gainers": {
      "get": {"summary": "Top gainers", "responses": {"200": {"description": "OK"}}}
    },
    "/test/security/{symbol}/company": {
      "get": {
        "summary": "Company details by symbol",
        "parameters": [{"name":"symbol","in":"path","required":true,"schema":{"type":"string"}}],
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/test/security/{symbol}/depth": {
      "get": {
        "summary": "Market depth by symbol",
        "parameters": [{"name":"symbol","in":"path","required":true,"schema":{"type":"string"}}],
        "responses": {"200": {"description": "OK"}}
      }
    },
    "/test/security/{symbol}/history": {
      "get": {
        "summary": "Price history by symbol",
        "parameters": [
          {"name":"symbol","in":"path","required":true,"schema":{"type":"string"}},
          {"name":"start","in":"query","required":true,"schema":{"type":"string","format":"date"}},
          {"name":"end","in":"query","required":true,"schema":{"type":"string","format":"date"}}
        ],
        "responses": {"200": {"description": "OK"}}
      }
    }
  }
}`
