package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niroopreddym/HLSVideoStreaming/handlers"
)

func main() {
	router := mux.NewRouter()
	cacheHandler := handlers.NewVideoCacheHandler(router)

	fmt.Println("HLS Stream Server : ", 9293)

	router.Handle("/asset/video/{id}/video", http.HandlerFunc(cacheHandler.GetPlaylistInfo)).Methods("GET")
	router.Handle("/asset/video/{id}/{segment_id}", http.HandlerFunc(cacheHandler.StreamFileSegments)).Methods("GET")

	fmt.Printf("Starting server on %v\n", 9293)

	// serve and log errors
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", 9293), router))
}

func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}
