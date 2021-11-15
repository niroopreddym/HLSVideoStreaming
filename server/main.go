package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niroopreddym/HLSConversion/handlers"
)

func main() {
	router := mux.NewRouter()
	cacheHandler := handlers.NewVideoCacheHandler()

	fmt.Println("HLS Stream Server : ", 9293)

	router.Handle("/asset/video/{id}/video.m3u8", http.HandlerFunc(cacheHandler.GetPlaylistInfo)).Methods("GET")
	router.Handle("/asset/video/{id}/{segment_id}", http.HandlerFunc(cacheHandler.StreamFileSegments)).Methods("GET")

	http.ListenAndServe(":9293", router)
}
