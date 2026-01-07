package handler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
)

type DebugRequest struct {
	Method        string
	URL           string
	Proto         string
	Header        http.Header
	Host          string
	RemoteAddr    string
	RequestURI    string
	ContentLength int64
	TransferEnc   []string
	Close         bool
	Query         url.Values
	Form          url.Values
	PostForm      url.Values
	TLS           *tls.ConnectionState
}

var router = setupRouter()

func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}

func setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /ip", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			w.Write(getIP("https://api.ipify.org"))
			return
		}

		w.Write(getIP(url))
	})

	router.HandleFunc("GET /hw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	router.HandleFunc("GET /debug", func(w http.ResponseWriter, r *http.Request) {
		d := DebugRequest{
			Method:        r.Method,
			URL:           r.URL.String(),
			Proto:         r.Proto,
			Header:        r.Header,
			Host:          r.Host,
			RemoteAddr:    r.RemoteAddr,
			RequestURI:    r.RequestURI,
			ContentLength: r.ContentLength,
			TransferEnc:   r.TransferEncoding,
			Close:         r.Close,
			Query:         r.URL.Query(),
			Form:          r.Form,
			PostForm:      r.PostForm,
			TLS:           r.TLS,
		}

		b, _ := json.MarshalIndent(d, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	router.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("https://pixeldrain.com/api/file/%s", id), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		req.Header.Set("Host", "pixeldrain.com")
		req.Header.Set("Origin", "https://pixeldrain.com")
		req.Header.Set("Referer", "https://pixeldrain.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/237.84.2.178 Safari/537.36")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		maps.Copy(w.Header(), res.Header)
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})

	return router
}

func getIP(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return body
}
