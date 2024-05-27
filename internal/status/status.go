package status

import (
	"log"
	"net/http"
)

func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	// an example API handler
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}
