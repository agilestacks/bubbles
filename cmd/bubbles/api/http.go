package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/agilestacks/bubbles/cmd/bubbles/config"
)

type middleware func(http.Handler) http.Handler

func mw(mws ...middleware) middleware {
	return func(handler http.Handler) http.Handler {
		h := handler
		for i := len(mws) - 1; i >= 0; i-- {
			h = mws[i](h)
		}
		return h
	}
}

func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		httptest.NewRecorder()

		if config.Debug {
			log.Printf("HTTP <<< %s %s", req.Method, req.URL)
		}
		crw := NewCapturingResponseWriter(rw, false)
		handler.ServeHTTP(crw, req)
		if config.Debug {
			log.Printf("HTTP === %d", crw.Captured.Status)
		}
	})
}

func withApiSecret(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !checkApiSecret(req) {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(rw, req)
	})
}

func withAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !checkApiSecret(req) {
			rw.Header().Set("WWW-Authenticate", "Basic realm=\".\"")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(rw, req)
	})
}

func allowedOrigins() []string {
	return []string{"*"}
}

func allowedHeaders() []string {
	return []string{
		"Accept", "Accept-Encoding", "Accept-Language", "Authorization",
		"Connection", "Content-Length", "Content-Type", "Token", "Session", "Host", "Origin",
		"X-CSRF-Token", "X-Requested-With", "X-Agent-Request", "X-Agent",
	}
}

func withAccessControl(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins(), ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders(), ","))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func getRouter() http.Handler {
	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r), withAccessControl)
	r.NotFoundHandler = mw(withLogger)(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusNotFound)
	}))

	cmw := mw(withLogger, withApiSecret)

	r.Handle("/", cmw(http.HandlerFunc(newBubble))).
		Methods(http.MethodPost, http.MethodOptions)

	s := r.PathPrefix("/{name}").Subrouter()
	s.Handle("", cmw(http.HandlerFunc(createBubble))).
		Methods(http.MethodPut, http.MethodOptions)
	s.Handle("", mw(withLogger)(http.HandlerFunc(getBubble))).
		Methods(http.MethodGet, http.MethodOptions)
	// s.Handle("", cmw(http.HandlerFunc(deleteBubble))).
	// 	Methods("DELETE")

	s = r.PathPrefix("/api/v1/ping").Subrouter()
	s.Handle("", mw(withLogger)(http.HandlerFunc(ping))).
		Methods(http.MethodGet, http.MethodOptions)

	return r
}

func ping(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func Listen(host string, port int) {
	r := getRouter()

	http.Handle("/", r)

	go listen(&http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	})
}

func listen(server *http.Server) {
	log.Fatalf("Unable to create HTTP server: %v", server.ListenAndServe())
}

func writeError(w http.ResponseWriter, status int, message string) {
	if config.Debug {
		log.Printf("Error %d HTTP: %s", status, message)
	}

	b, err := json.Marshal(struct {
		Error string `json:"error"`
	}{message})

	if err != nil {
		msg := fmt.Sprintf("Unable to marshall JSON: %v", err)
		log.Print(msg)
		b = []byte(msg)
		w.Header().Set("Content-Type", "text/plain")
		status = http.StatusInternalServerError
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.WriteHeader(status)
	w.Write(b)
}

func checkApiSecret(req *http.Request) bool {
	if config.BubblesApiSecret == "" {
		return true
	}

	xApiSecret := req.Header.Get("X-API-Secret")
	if config.Trace {
		log.Printf("X-API-Secret `%v`", xApiSecret)
	}
	if xApiSecret == config.BubblesApiSecret {
		return true
	}

	username, password, ok := req.BasicAuth()
	if config.Trace {
		log.Printf("Authorization `%v`; Basic auth %v `%v` `%v`",
			req.Header.Get("Authorization"), ok, username, password)
	}
	if ok {
		return username == config.BubblesApiSecret || password == config.BubblesApiSecret
	}

	return false
}
