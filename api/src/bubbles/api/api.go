package api

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	// "bubbles/config"
)

func createBubble(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	id, exist := vars["name"]
	if !exist || id == "" {
		if gen, ok := idOrError(w); !ok {
			return
		} else {
			id = gen
		}
	}
	ttl := 10 // TODO ?ttl=
	create(w, req, id, ttl)
}

func newBubble(w http.ResponseWriter, req *http.Request) {
	id, ok := idOrError(w)
	if !ok {
		return
	}
	ttl := 10
	create(w, req, id, ttl)
}

func idOrError(w http.ResponseWriter) (string, bool) {
	id, err := RandomString(10)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return "", false
	}
	return id, true
}

func create(w http.ResponseWriter, req *http.Request, id string, ttl int) {
	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		writeError(w, http.StatusInternalServerError,
			fmt.Sprintf("Error reading request body: %v", err))
		return
	}

	if len(body) == 0 {
		log.Printf("WARN: storing an empty `%s`", id)
	}

	err = store(id, body, ttl)

	if err != nil {
		message := fmt.Sprintf("Unable to store `%s`: %v", id, err)
		log.Print(message)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "exist") {
			status = http.StatusConflict
		} else if strings.Contains(err.Error(), "supported") {
			status = http.StatusBadRequest
		}
		writeError(w, status, message)
	} else {
		w.Header().Add("Location", "/"+id)
		w.WriteHeader(http.StatusCreated)
	}
}

func getBubble(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	id, exist := vars["name"]
	if !exist || id == "" {
		writeError(w, http.StatusBadRequest, "No bubble id specified")
	}

	body, err := retrieve(id)
	if err != nil {
		message := fmt.Sprintf("Unable to retrieve `%s`: %v", id, err)
		log.Print(message)
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "found") {
			status = http.StatusNotFound
		}
		writeError(w, status, message)
	} else {
		if body == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		wrote, err := w.Write(body)
		if err != nil || wrote != len(body) {
			log.Printf("Unable to send `%s` (wrote %d out of %d bytes): %v",
				wrote, len(body), err)
		}
	}
}
