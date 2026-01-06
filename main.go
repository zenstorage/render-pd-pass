package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
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

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /ip", func(w http.ResponseWriter, r *http.Request) {
		ip := getIP()
		w.Write(ip)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, router)
}

func getIP() []byte {
	res, err := http.Get("https://api.ipify.org")
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
